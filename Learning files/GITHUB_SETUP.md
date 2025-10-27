# GitHub Setup Instructions

## Step 1: Authenticate with GitHub

Run this command in your terminal:
```bash
gh auth login
```

When prompted:
1. Select **"GitHub.com"** (press Enter)
2. Choose **"HTTPS"** (press Enter)
3. Choose **"Yes"** to authenticate Git with your GitHub credentials
4. Choose **"Login with a web browser"** (most reliable option)
5. Copy the code that appears
6. Press Enter to open your browser
7. Paste the code in the browser and authorize
8. Return to terminal - you should see "Successfully authenticated!"

## Step 2: Create Repository and Push

After authentication, run these commands:

```bash
# This creates the repo on GitHub and pushes your code
gh repo create go-todo-api --public --source=. --push
```

That's it! Your code will be on GitHub!

## Alternative: Manual Setup

If you prefer manual setup:

1. Create repo on GitHub.com: https://github.com/new
   - Name: `go-todo-api`
   - Choose Public or Private
   - Don't initialize with README
   - Click "Create repository"

2. Run these commands:
```bash
git remote add origin https://github.com/YOUR_USERNAME/go-todo-api.git
git branch -M main
git push -u origin main
```

