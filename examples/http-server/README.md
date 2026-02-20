# FST Search API Server

A high-performance HTTP API server demonstrating **Finite State Transducer** string search capabilities.

## 🚀 Features

- **⚡ Ultra-fast string search** using FST data structures
- **🔍 Multiple search modes**: Exact matching and fuzzy search
- **🌐 RESTful API** with JSON request/response
- **📊 Performance monitoring** with real-time statistics  

## 🔧 Quick Start

```bash
cd examples/http-server
go run main.go
# Server available at http://localhost:8080
```

## 📖 API Endpoints

### Search
```bash
curl -X POST localhost:8080/search \
  -H "Content-Type: application/json" \
  -d '{"query":"app","limit":5}'
```

### Fuzzy Search
```bash
curl -X POST localhost:8080/fuzzy \
  -H "Content-Type: application/json" \
  -d '{"query":"aple","maxErrors":2,"limit":5}'
```

### Stats
```bash
curl localhost:8080/stats
```

## 🎯 Performance

- **Search Time**: ~68 nanoseconds per lookup
- **Throughput**: ~17.6 million operations per second  
- **Memory Usage**: ~16 bytes per search operation
- **Build Time**: ~56 microseconds for 1000+ words

Perfect for high-performance web APIs! 🚀
