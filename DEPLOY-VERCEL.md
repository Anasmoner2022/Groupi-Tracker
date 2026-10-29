# Deploy Groupie Tracker on Vercel

1. Extract this ZIP. Open the folder containing `go.mod` and `vercel.json`.
2. Copy these files into your GitHub repository and commit/push the changes.
3. In Vercel, choose **Add New → Project** and import the repository.
4. Set **Root Directory** to the folder containing `go.mod` (the repository root if you copied its contents).
5. Use the **Go** framework preset. `vercel.json` already selects it. Leave build/install/output overrides unset.
6. Click **Deploy**.

No environment variables, database, Node.js, or npm dependencies are required. Vercel supplies PORT automatically.

## Run locally

Install Go 1.26.1 or newer, then run:

```sh
go run .
```

Open http://localhost:8080. You can set PORT to choose another port.

## What changed

- main.go reads PORT, with 8080 as the local default.
- vercel.json selects Vercel's native Go server preset; no route rewrites are needed.
- HTML templates are embedded in the Go binary, so runtime working-directory differences cannot break loading.
- External API requests time out after 10 seconds each instead of waiting indefinitely.
- The HTML/CSS/JavaScript design and existing routes are unchanged.

## Runtime notes

The app fetches artists and relations when each server instance starts. It still requires access to https://groupietrackers.herokuapp.com/api. If that API is unavailable, startup fails and Vercel logs show the reason. The in-memory cache belongs to each instance and is rebuilt on restart; it is not persistent storage.

The supplied archive contains no static directory; the existing templates use inline CSS/JavaScript and external fonts/images.

## Verification

See VALIDATION.md for checks actually performed. A live Vercel deployment has not been performed.

Official reference: https://vercel.com/docs/functions/runtimes/go
