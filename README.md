# check24-GenDev-7: Resilient Internet Provider Comparison

[![Go Server Coverage](https://github.com/rotmanjanez/check24-gendev-7/blob/badges/server/coverage-go.svg)](https://github.com/rotmanjanez/check24-gendev-7/tree/badges)
[![🚀 CI/CD Go Server](https://github.com/rotmanjanez/check24-gendev-7/actions/workflows/ci-cd-go-server.yml/badge.svg)](https://github.com/rotmanjanez/check24-gendev-7/actions/workflows/ci-cd-go-server.yml)
[![🚀 CI/CD Client Website](https://github.com/rotmanjanez/check24-gendev-7/actions/workflows/ci-cd-client-website.yml/badge.svg)](https://github.com/rotmanjanez/check24-gendev-7/actions/workflows/ci-cd-client-website.yml)

---
## Challenge Accepted!

This project is my submission for the **7th round of the CHECK24 GenDev Internet Provider Comparison Challenge**. The core task was to build an application that allows users to seamlessly compare offers from five different internet providers, even when their APIs are unreliable or slow. The aim was to ensure a **smooth user experience** and **broad product coverage**, showing only bookable offers without sacrificing responsiveness.

---

## Implemented Features

| Feature                               | Notes                                                             |
| ------------------------------------- | ----------------------------------------------------------------- |
| **Robust API Failure/Delay Handling** | Parallel fetching, timeouts, retries, circuit breakers.           |
| **Sorting & Filtering**               | Filter by speed, price, duration; sort by various criteria.       |
| **Shareable Result Links**            | Short URLs persist offer states for consistent sharing.           |
| **API Credential Security**           | Server-side only, never logged above DEBUG level.                 |
| **User Input Validation**             | Form validation and clear feedback.                               |
| **Session State**                     | Stores user preferences and address using `sessionStorage`.       |
| **Personalization**                   | i18n (English and German), Dark Mode support.                     |
| **Mobile-Friendly UI**                | Responsive layout for small screens.                              |
| **Progressive Results Loading**       | Early offers shown while others load in the background.           |
| **Scalable & Resilient Backend**      | Stateless design using Redis enables scaling and smooth rollouts. |
| **Data Safety Measures**              | Limits on input and offer sizes prevent performance issues.       |

---

## Architecture & Tech Stack

The system is designed for maintainability and reliability using a carefully chosen stack:

**Backend (Go):**

* **Language:** Go (Golang) for performance and concurrency.
* **Routing:** Generated from OpenAPI (`openapi.yaml`) using `gorilla/mux`.
* **Caching:** Redis for persistence and shareable link support.
* **Resilience:** Custom request manager handles retries, fallbacks, and timeouts.

**Frontend (Vue):**

* **Framework:** Vue 3 with TypeScript.
* **Tooling:** Vite for fast builds.
* **Styling:** TailwindCSS.
* **UI Components:** Shadcn/ui.

**API Definition:**

* OpenAPI used to define a strict frontend-backend contract.
* Code generation ensures consistency.

**CI/CD & Deployment:**

* GitHub Actions for CI.
* Docker Compose for local use.
* Architecture ready for Kubernetes/Helm.

---

## API Design Philosophy

The API favors clarity, progressive delivery, and robustness:

### REST Interface

1. **Start Search:** `POST /internet-products` launches the provider search.
2. **Continue Fetching:** `GET /internet-products/continue` uses cursors to fetch progressive results.
3. **Share Results:** `POST /internet-products/share/{cursor}` saves a snapshot of results, accessible via a short link.

Note: The house number is an optional string, as e.g. `6a` is a valid house number and there are addresses without house number (e.g. `Pariser Platz, 10117 Berlin`).

### Backend Scalability & Resilience

The backend is designed as a **stateless** service that can be **horizontally scaled** with ease. Any instance can handle any request, with synchronization handled via a **Redis cache**, which also manages persistence. This setup supports various **load balancers** that don't require domain-specific knowledge to distribute traffic effectively. It also simplifies **slow rollouts**: no extra dependencies are needed, and newer backend versions can work with the same database and cache as the current deployment. Requests can be gradually routed to the updated version, allowing smooth transitions without complex migrations or service interruptions, and enabling easy rollback if needed.

Thanks to this design, even user queries that are already in flight could be migrated to another instance. This is only possible because of the internal interface used by providers:

### Provider Integration Interface

All providers implement a simple interface to convert:

* user queries into HTTP requests, and
* HTTP responses into internet products and follow-up responses.

The runtime centrally handles retries, error handling, and caching.

```go
type ProviderAdapter interface {
    PrepareRequest(ctx context.Context, request Request) (ParsedResponse, error)
    ParseResponse(ctx context.Context, response Response) (ParsedResponse, error)
    Name() string
}
```

This approach ensures that provider implementations remain focused purely on API semantics, resulting in **truly stateless** provider code. It significantly reduces the **implementation and maintenance** burden for both new and existing providers. Runtime improvements automatically benefit all providers.

For any given query, all state is simply a list of remaining requests, which makes tracing easy and debugging simple.

Note: This is the second key part to make mid-flight server migrations in the future straightforward: All state in a single place and simple datastructure with clear semantics.

### Developer Experience

Go provides a great developer experience, with fast build times and strong debugging capabilities. The project's documentation is centered around the most important parts: the provider interface and the request manager. In other areas, Go's simplicity makes the implementation clear enough that **function names often speak for themselves**, which is why I chose not to extensively document every single function. The code should be simple and self-explanatory.

In addition, **E2E, unit, and integration tests** bolster product stability and give developers confidence when making changes. Test coverage isn’t a perfect 90%+, but it’s focused where it counts: provider-specific logic, edge cases, and end-to-end behavior using **mock providers**.


---

## Noteworthy Aspects

* **Mobile-Ready UI:** Minimal, responsive design ensures usability across devices.
* **Resource Efficient:** The current deployment is working great on an old 1 core X86 CPU with 1GB Ram (Oracle Cloud Free Tier Server).
* **Deliberate Backend Dependencies:** Only critical packages (`go-redis`, `gorilla/mux`) are used, minimizing third-party risk.

**<\ModernWebRant>**
*For the client, I really wish it was reasonably possible to have a decent modern experience without *so* many dependencies. Fun fact (or maybe not so fun): my package-lock.json is pushing over 6000 lines – that's like a quarter of my whole codebase!*
**<\/ModernWebRant>**

---

## Getting Started

### Prerequisites

* Go
* Node.js / npm 
* Docker & Docker Compose (if you are into that)

### Local Setup

Clone the repositorey
```bash
git clone git@github.com:rotmanjanez/check24-gendev-7.git
```
<table>
  <thead>
    <tr>
      <th>Docker</th>
      <th>Manual</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>
        <pre><code>cd check24-gendev-7
docker-compose up --build -d
</code></pre>
      </td>
      <td>
        <pre><code>cd server
go mod tidy
go run cmd/check24-gendev-7-server/main.go -localdev
</code></pre>
        <pre><code>cd client
npm install
npm run dev
</code></pre>
      </td>
    </tr>
  </tbody>
</table>


## Configuration

* **Server Settings:** `server/config.json` for timeouts, cache settings, provider configs.
* **API Keys:** Stored in `.env` file, never committed or logged.

---

## Future Work

* **URL Encode Aktive filters in share link:*
  Improve customer experience by storing current filter settings in the share link.

* **Address Autocomplete:**
  A custom fuzzy search using OpenStreetMap data could replace reliance on costly, branded third-party APIs like Google Places.

* **AI-Based Filtering:**
  A user intent parser using natural language processing could let users search with phrases like "streaming and gaming for four people," translating them into filter parameters. Wether this is worth from a buisness perspective is up to debate, as inference costs might be too high for the product usecase.

---

## Credits

Developed by **Janez Rotman** for the CHECK24 GenDev-7 Challenge.
Contact: [janez@janez.at](mailto:janez@janez.at)
