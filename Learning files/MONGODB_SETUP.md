# MongoDB Setup Guide

## Option 1: Use MongoDB Atlas (Cloud - Easiest) ⭐ Recommended

MongoDB Atlas is free and doesn't require local installation!

### Steps:
1. Go to https://www.mongodb.com/cloud/atlas/register
2. Sign up for free
3. Create a free cluster (M0)
4. Create a database user
5. Add your IP to the whitelist (or 0.0.0.0/0 for testing)
6. Get your connection string
7. Use that string in our code

### Connection String Format:
```
mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
```

## Option 2: Install MongoDB Locally

### macOS (using Homebrew):
```bash
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

### Verify it's running:
```bash
mongosh
```

## What We'll Change in the Code

### Current (In-Memory):
- Tasks stored in a `[]Task` slice
- Lost when server restarts
- Simple for learning

### With MongoDB:
- Tasks stored in MongoDB database
- Persists between restarts
- Production-ready

### Key Changes:
1. **Import MongoDB packages**
2. **Connect to MongoDB** in `main()`
3. **Use MongoDB IDs** instead of integers
4. **Read/write from database** in each handler

## Testing Without MongoDB First

You can test our current in-memory version and add MongoDB later when you're ready!

Just run:
```bash
go run main.go
```

