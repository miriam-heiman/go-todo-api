# AWS Lambda Deployment Guide

## Overview

This API can run in two modes:
1. **HTTP Server** (Docker) - For local development
2. **AWS Lambda** (Serverless) - For production deployment

## Prerequisites

### 1. Install Serverless Framework
```bash
npm install -g serverless
```

### 2. Configure AWS Credentials
```bash
# Option A: AWS CLI
aws configure

# Option B: Environment variables
export AWS_ACCESS_KEY_ID=your_key_here
export AWS_SECRET_ACCESS_KEY=your_secret_here
```

### 3. Set Up MongoDB Atlas

1. Create account at [mongodb.com/atlas](https://www.mongodb.com/atlas)
2. Create a new cluster (free M0 tier available)
3. Create database user
4. Whitelist IP: `0.0.0.0/0` (allow from anywhere - Lambda IPs change)
5. Get connection string:
   ```
   mongodb+srv://username:password@cluster.mongodb.net/todoapi?retryWrites=true&w=majority
   ```

### 4. Configure Environment Variables

Create `.env.production` file:
```bash
MONGO_URI=mongodb+srv://user:pass@cluster.mongodb.net/todoapi
API_KEY=your-production-api-key
OTEL_EXPORTER_OTLP_ENDPOINT=https://tempo.yourcompany.com:4318
LOKI_ENDPOINT=https://loki.yourcompany.com:3100
```

## Building for Lambda

### Build ARM64 (Recommended - Cheaper)
```bash
make build-lambda
```

### Build AMD64 (Intel/AMD)
```bash
make build-lambda-amd64
```

This creates a `bootstrap` binary that Lambda will execute.

## Deployment

### Deploy to Dev Environment
```bash
# Export environment variables
export MONGO_URI="mongodb+srv://..."
export API_KEY="your-api-key"

# Deploy
make deploy-lambda-dev
```

### Deploy to Production
```bash
# Export environment variables from production .env
export $(cat .env.production | xargs)

# Deploy
make deploy-lambda-prod
```

### Manual Deployment
```bash
# Build
make build-lambda

# Deploy
serverless deploy --stage prod --region us-east-1
```

## Post-Deployment

### Get API Endpoint
After deployment, Serverless will output:
```
endpoints:
  ANY - https://abc123.execute-api.us-east-1.amazonaws.com/dev/{proxy+}
  ANY - https://abc123.execute-api.us-east-1.amazonaws.com/dev
```

### Test the API
```bash
# Health check
curl https://your-api-url.amazonaws.com/dev/health

# Get tasks (requires API key)
curl -H "X-API-Key: your-api-key" https://your-api-url.amazonaws.com/dev/tasks

# Create task
curl -X POST \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test from Lambda","description":"Testing"}' \
  https://your-api-url.amazonaws.com/dev/tasks
```

## Monitoring & Observability

### CloudWatch Logs
```bash
# View Lambda logs
serverless logs -f api --tail

# View specific time range
serverless logs -f api --startTime 1h
```

### AWS X-Ray (Distributed Tracing)
- Lambda is configured with X-Ray enabled
- View traces in AWS Console: X-Ray → Traces
- Trace ID appears in CloudWatch logs

### Connect Grafana to AWS
Configure Grafana to query:
- **CloudWatch Logs** (replaces Loki for Lambda logs)
- **X-Ray** (replaces Tempo for Lambda traces)

**Grafana CloudWatch Setup:**
1. Add CloudWatch data source
2. Configure AWS credentials
3. Query logs: `/aws/lambda/go-todo-api-dev-api`

## Cost Estimation

### Free Tier
- 1M requests/month free
- 400,000 GB-seconds compute free

### Beyond Free Tier (ARM64)
- **Requests:** $0.20 per 1M requests
- **Compute:** $0.0000133334 per GB-second
- **Example:** 10M requests/month with 512MB, 200ms avg:
  - Requests: 10M × $0.20 = $2.00
  - Compute: 10M × 0.5GB × 0.2s × $0.0000133334 = $13.33
  - **Total: ~$15/month**

Compare to EC2: t3.micro = $8.50/month (but can't auto-scale)

## Troubleshooting

### Lambda Cold Starts
**Problem:** First request after idle is slow (200-500ms)

**Solutions:**
1. Use ARM64 (faster cold starts)
2. Keep Lambda warm with scheduled pings
3. Increase memory (1024MB = faster CPU)

### MongoDB Connection Issues
**Problem:** `connection timeout` or `no reachable servers`

**Solutions:**
1. Whitelist `0.0.0.0/0` in MongoDB Atlas
2. Use connection string from Atlas dashboard
3. Check VPC settings if using private MongoDB

### "Permission Denied" Errors
**Problem:** `AccessDeniedException` in logs

**Solution:** Update IAM role in `serverless.yml`:
```yaml
provider:
  iam:
    role:
      statements:
        - Effect: Allow
          Action:
            - your-required-action
          Resource: '*'
```

### High Costs
**Problem:** Unexpected AWS bill

**Solutions:**
1. Check CloudWatch Logs retention (default: never expire)
2. Set log retention: 7 days for dev, 30 days for prod
3. Monitor with AWS Cost Explorer
4. Set billing alerts

## Cleanup

### Remove Deployment
```bash
make remove-lambda
```

Or:
```bash
serverless remove --stage dev
```

This removes:
- Lambda function
- API Gateway
- CloudWatch log groups
- IAM roles

## CI/CD Integration

### GitHub Actions
The CI/CD pipeline (`.github/workflows/ci.yml`) can be extended for Lambda deployment:

```yaml
deploy-lambda:
  name: Deploy to AWS Lambda
  runs-on: ubuntu-latest
  needs: test
  if: github.ref == 'refs/heads/main'

  steps:
    - uses: actions/checkout@v4

    - uses: actions/setup-go@v5
      with:
        go-version: '1.23'

    - name: Build Lambda
      run: make build-lambda

    - name: Configure AWS
      uses: aws-actions/configure-aws-credentials@v4
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: us-east-1

    - name: Deploy
      env:
        MONGO_URI: ${{ secrets.MONGO_URI }}
        API_KEY: ${{ secrets.API_KEY }}
      run: serverless deploy --stage prod
```

## Architecture Diagram

```
User Request
    ↓
AWS API Gateway (HTTPS endpoint)
    ↓
AWS Lambda (Go function)
    ├─→ MongoDB Atlas (Database)
    ├─→ Tempo (Traces via OTLP)
    └─→ Loki (Logs via HTTP)
         ↓
    CloudWatch Logs (Backup)
    CloudWatch X-Ray (Backup traces)
         ↓
    Grafana (Visualization)
```

## Comparison: Docker vs Lambda

| Feature | Docker (Local) | Lambda (Production) |
|---------|---------------|---------------------|
| **Cost** | $0 (local) | ~$15/month (10M requests) |
| **Scaling** | Manual | Automatic (0-1000s instances) |
| **Maintenance** | You manage | AWS manages |
| **Cold Start** | No | Yes (~200ms) |
| **Max Timeout** | Unlimited | 15 minutes |
| **Observability** | Loki+Tempo (local) | CloudWatch+X-Ray (or Loki+Tempo remote) |
| **Best For** | Development | Production, variable traffic |

## Next Steps

1. **Set up MongoDB Atlas** (free tier)
2. **Deploy to dev** with `make deploy-lambda-dev`
3. **Test the endpoint**
4. **Configure Grafana** to connect to CloudWatch
5. **Set up CI/CD** for automatic deployments
6. **Deploy to prod** when ready

## Support

For issues or questions:
- Check CloudWatch logs: `serverless logs -f api --tail`
- Review X-Ray traces in AWS Console
- Check MongoDB Atlas metrics
- Test locally: `make run-local`
