# angkot-ts

To install dependencies:

```bash
bun install
```

To run:

```bash
bun run index.ts
```

This project was created using `bun init` in bun v1.1.30. [Bun](https://bun.sh) is a fast all-in-one JavaScript runtime.

To build:

```bash
docker build -t angkot-bot .
```

To run in docker:

```bash
docker run -e TELEGRAM_TOKEN="your_telegram_token" -e ADMIN_ID="your_admin_id" -v $(pwd)/database:/app/database --restart unless-stopped angkot-bot
```
