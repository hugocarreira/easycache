# 🛠️ Contributing to EasyCache

Thank you for your interest in contributing to EasyCache! 🎉  
We welcome all contributions, whether it's bug reports, feature requests, or pull requests.  

---

## 📌 Reporting Issues

If you found a bug or have a question, please open an **[Issue](https://github.com/hugocarreira/easycache/issues)**.  
Make sure to include:  
✅ A **clear description** of the problem.  
✅ Steps to **reproduce the issue** (if applicable).  
✅ Your **Go version** and **system details**.  

Before opening a new issue, **check if it hasn't been reported** already.  

---

## 💡 Requesting a Feature

We love new ideas! 🚀 If you have a suggestion, open an **Issue** with:  
✅ A **detailed explanation** of the feature.  
✅ Why this feature is useful.  
✅ Example use cases.  

## 🔄 Submitting a Pull Request (PR)

✅ Open your pull request against `main`.  
✅ Create a Pull Request with a clear description of your changes.


### 🧪 Running Tests & Benchmarks

Before submitting code, ensure all tests pass:

```sh
go test ./... -v
```

```sh
go test -race ./...
```

```sh
go vet ./...
```

```sh
go test -coverpkg=./... ./tests
```

####  🚀 Performance Benchmarks  

Run benchmarks locally when evaluating performance-sensitive changes:

```sh
go test -run '^$' -bench=. -benchmem ./tests
```

## 🎯 Contributing Best Practices

✅ Be respectful – We welcome contributions from everyone!  
✅ Keep discussions focused – Avoid unrelated topics in issues/PRs.  
✅ Improve documentation – Even small fixes help!  
✅ Test your code – Ensure everything works before submitting.  
