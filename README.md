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