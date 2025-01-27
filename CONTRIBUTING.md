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

✅ Open your pull request against `master`.  
✅ Create a Pull Request with a clear description of your changes.


### 🧪 Running Tests & Benchmarks

Before submitting code, ensure all tests pass:

```sh
go test ./tests -v
```

```sh
go test -bench=. -benchmem ./tests
```

####  🚀 Performance Benchmarks

We ran performance benchmarks on EasyCache to measure the efficiency of `Set()`, `Get()`, `Delete()`, and eviction policies (`FIFO`, `LRU`, `LFU`).

| Benchmark                | Iterations  | Time per operation | Memory used | Allocations per op |
|--------------------------|------------|--------------------|-------------|--------------------|
| **`BenchmarkCacheSet`**   | 2,936,356  | **408.4 ns/op**    | **122 B/op**  | **5 allocs/op** |
| **`BenchmarkCacheGet`**   | 39,143,538 | **30.79 ns/op**    | **0 B/op**    | **0 allocs/op** |
| **`BenchmarkCacheDelete`**| 5,376,940  | **223.3 ns/op**    | **96 B/op**   | **3 allocs/op** |
| **`BenchmarkFIFOEviction`** | 3,065,480 | **391.7 ns/op**    | **122 B/op**  | **5 allocs/op** |
| **`BenchmarkLRUEviction`** | 3,045,759  | **402.1 ns/op**    | **122 B/op**  | **5 allocs/op** |
| **`BenchmarkLFUEviction`** | 2,916,150  | **394.3 ns/op**    | **88 B/op**   | **4 allocs/op** |

**Tested on:**  
- **Go Version:** 1.23.5  
- **Cache Configuration:** `MaxSize = 10,000`, `TTL = 60s`  

---

##### 🛠️ Running the Benchmarks  

To run the benchmarks yourself, use:  

```sh
go test -bench=. -benchmem ./tests
```

## 🎯 Contributing Best Practices

✅ Be respectful – We welcome contributions from everyone!  
✅ Keep discussions focused – Avoid unrelated topics in issues/PRs.  
✅ Improve documentation – Even small fixes help!  
✅ Test your code – Ensure everything works before submitting.  

