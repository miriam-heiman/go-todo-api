# MongoDB Quick Start Guide

## Get Your Connection String

### Step 1: Sign up for MongoDB Atlas (Free)
1. Go to: https://www.mongodb.com/cloud/atlas/register
2. Sign up with your email (it's completely free!)
3. Choose the **FREE M0** cluster option

### Step 2: Create a Database User
1. Click "Database Access" in left sidebar
2. Click "Add New Database User"
3. Set username and password (remember these!)
4. Click "Add User"

### Step 3: Whitelist Your IP
1. Click "Network Access" in left sidebar
2. Click "Add IP Address"
3. Click "Allow Access from Anywhere" (for testing)
4. Click "Confirm"

### Step 4: Get Your Connection String
1. Click "Database" in left sidebar
2. Click "Connect" on your cluster
3. Choose "Connect your application"
4. Copy the connection string (looks like this):

```
mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
```

### Step 5: Update Your Code
Replace `YOUR_CONNECTION_STRING` in main.go with your actual connection string!

---

## Quick Commands

**Once you have the connection string, update line 45 in main.go**

Run your server:
```bash
go run main.go
```

Your data will now persist in the cloud! 🎉

