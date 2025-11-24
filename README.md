# Study-Llama

Study-Llama is a demo application for organizing, extracting, and searching study notes using LlamaAgents, LlamaClassify, and LlamaExtract. It features a Go-based web frontend and a Python backend for advanced note processing.

## Overview

- **Frontend (Go):**  
  Serves the web interface, handles user authentication, file uploads, and search requests. Renders pages using HTML templates and manages static assets.
- **Backend (Python):**  
  Provides APIs for classifying notes, extracting information, storing/retrieving data, and performing metadata-filtered vector searches. Handles database operations and vector search.

## How It Works

1. **User Interaction:**  
   Users access the web UI to sign in, create categories for their study notes, upload notes, and search among them. 

2. **Frontend-to-Backend Communication:**  
   The Go server receives requests from the UI and communicates with the Python backend via HTTP API calls.  
   - For example, when a user uploads a note, the frontend sends the file and metadata to the backend for processing.
   - Search queries are also forwarded to the python backend for vector search.

3. **Backend Processing:**  
   The Python backend (deployed on LlamaCloud as a LlamaAgent) handles:
   - **Classification & Extraction:**  
     Uses workflows in `classify_and_extract/workflow.py` to classify notes (with LlamaClassify) and extract structured information (with LlamaExtract).
   - **Database Operations:**  
     - **Files:**  
       Uploaded files and their metadata (name, category, owner) are stored.
     - **Classification Rules:**  
       Custom rules for categorizing notes are stored and retrieved from the database.
   - **Vector Search:**  
     Extracted summaries and FAQs are indexed for semantic search and retrieval.

## Features

- Upload and categorize study notes.
- Extract structured information from notes.
- Search notes with metadata filters.
- User authentication and access control.
- Modern web UI with Go templates.

## Project Structure

- **frontend/**  
  - `main.go`: Web server entry point.  
  - `handlers/`: HTTP request handlers.  
  - `auth/`, `files/`, `rules/`: Business logic and DB operations.  
  - `static/`: Images and assets.  
  - `templates/`: HTML templates.
- **src/**  
  - `study_llama/`: Python backend modules for classification, extraction, search, and database logic.

## Set Up

**Clone this repository:**

```bash
git clone https://github.com/run-llama/study-llama
cd /study-llama
```

**Deploy the LlamaAgent:**

In order for the classify-extract and vector search agent workflows to receive and process queries, they needs to be deployed to the cloud (or at least accessible through a public endpoint). The easiest way to do so is to use [`llamactl`](https://developers.llamaindex.ai/python/llamaagents/llamactl/getting-started/) and deploy the agent workflow as a [LlamaAgent](https://developers.llamaindex.ai/python/llamaagents/overview/):

```bash
uv tool install -U llamactl
llamactl auth # authenticate
llamactl deployments create # create a deployment from the current repository
```

In order for the LlamaAgent to work, you will need the following environment variables in a `.env` file (`llamactl` manages environments autonomously):

- `OPENAI_API_KEY` to interact with GPT-4.1 for email generation
- `LLAMA_CLOUD_API_KEY` and `LLAMA_CLOUD_PROJECT_ID` to get predictions from LlamaClassify and LlamaExtract
- `POSTGRES_CONNECTION_STRING` to connect to the Postgres database with the uploaded files and the classification rules (you can use [Neon](https://neon.com), [Supabase](https://supabase.com), [Prisma](https://prisma.io) or a self-hosted Postgres instance)
- `QDRANT_API_KEY` and `QDRANT_HOST`, to upload sumamries and question/answers to perform vector search and retrieval (you can use [Qdrant Cloud](https://qdrant.tech) or a self-hosted Qdrant instance).

**Deploy the frontend**

Once the agent is deployed, build the Docker image for the frontend (needed to interact with the LlamaAgents we just created), and deploy it through services like [Dokploy](https://dokploy.com) or [Coolify](https://coolify.io).

```bash
docker build . -t your-username/study-llama:prod
# docker login ghcr.io # (uncomment if you wish to use the GitHub container registry)
docker push your-username/study-llama:prod
```

The frontend service uses a few env variables:

- `LLAMA_CLOUD_API_KEY`, `FILES_API_ENDPOINT` (which will presumably be `https://api.cloud.llamaindex.ai/deployments/study-llama/workflows/classify-and-extract/run`) and `SEARCH_API_ENDPOINT` (which will presumably be `https://api.cloud.llamaindex.ai/deployments/study-llama/workflows/search/run`), the API key and the API endpoints to interact with your deployed LlamaAgent
- `POSTGRES_CONNECTION_STRING` to connect to the Postgres database with the uploaded files, the classification rules and the user auth (you can use [Neon](https://neon.com), [Supabase](https://supabase.com), [Prisma](https://prisma.io) or a self-hosted Postgres instance, but it has to be the **same as for the LlamaAgent**)
- `CACHE_TABLE` and `RATE_LIMITING_TABLE`, the table names for the SQLite database taking care of caching and rate limiting.

Services like Dokploy or Coolify offer you to set these environment variables through their own environment management interfaces.