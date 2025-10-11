# 📚 Go E-Commerce API - Student Documentation

> **Welcome!** This documentation is written for students coming from a Node.js background who want to learn Go by building a real-world e-commerce API.

## 📖 Documentation Index

This documentation is organized into easy-to-follow sections:

1. **[Getting Started](./01-getting-started.md)** - Setup and running your first API
2. **[Project Structure](./02-project-structure.md)** - Understanding the folder organization
3. **[Go vs Node.js](./03-go-vs-nodejs.md)** - Key differences explained
4. **[Database Setup](./04-database-setup.md)** - PostgreSQL, migrations, and models
5. **[Authentication](./05-authentication.md)** - JWT and user management
6. **[API Endpoints](./06-api-endpoints.md)** - All available routes
7. **[Code Walkthrough](./07-code-walkthrough.md)** - Deep dive into the codebase
8. **[Common Tasks](./08-common-tasks.md)** - How to add features
9. **[Testing & Debugging](./09-testing-debugging.md)** - Development tips
10. **[Deployment](./10-deployment.md)** - Going to production

## 🎯 Learning Path

### For Complete Beginners in Go:
1. Start with **Getting Started** (01)
2. Read **Go vs Node.js** (03) to understand the paradigm shift
3. Study **Project Structure** (02)
4. Follow **Database Setup** (04)
5. Explore **Code Walkthrough** (07)

### For Intermediate Developers:
1. Quick skim of **Getting Started** (01)
2. Read **Project Structure** (02)
3. Jump to **Code Walkthrough** (07)
4. Check **Common Tasks** (08) for practical examples

## 🚀 Quick Start

```bash
# 1. Clone and setup
git clone <repo-url>
cd goecom

# 2. Install dependencies
go mod download

# 3. Setup environment
cp .env.example .env
# Edit .env with your database URL

# 4. Run migrations
./scripts/migrate.sh up

# 5. Seed database
./scripts/seed.sh

# 6. Run the server
go run cmd/server/main.go
```

## 💡 What You'll Learn

- ✅ Building RESTful APIs with Go and Gin framework
- ✅ PostgreSQL database integration with GORM
- ✅ JWT authentication and middleware
- ✅ Project structure and organization
- ✅ Error handling and validation
- ✅ Database migrations
- ✅ Environment configuration
- ✅ HTTP routing and handlers

## 🎓 Prerequisites

- Basic programming knowledge (JavaScript/Node.js preferred)
- Understanding of REST APIs
- Familiarity with Git
- Basic command line usage

## 📝 Notes Convention

Throughout this documentation:
- 💡 **Tip** - Helpful hints and best practices
- ⚠️ **Warning** - Common pitfalls to avoid
- 🔍 **Compare** - Node.js vs Go comparison
- 📌 **Remember** - Key concepts to keep in mind
- 🎯 **Practice** - Hands-on exercises

## 🤝 Need Help?

- Check the **Common Tasks** section for practical examples
- Review **Testing & Debugging** for troubleshooting
- Look at the code comments - they're there to help!

Happy Learning! 🚀
