# 🚀 Instech --- Project Plan

## 🧠 Project Overview

StackAlt is a **developer decision tool** to explore, compare,
and choose programming languages, frameworks, and libraries.

Instead of being a static dataset, the product helps answer:

-   What should I use instead?
-   What is similar to this?
-   How are tools connected?

👉 The core concept: \> A **graph of developer tools + meaningful relationships**

------------------------------------------------------------------------

## 🎯 Core Value Proposition

-   Discover alternatives to tools
-   Understand ecosystems (React → Next.js, etc.)
-   Visualize relationships between technologies
-   Make better tech decisions faster

------------------------------------------------------------------------

## 🧩 Core Features

### Backend Features

-   Tool management (languages, frameworks, libraries)
-   Relationships management
-   Alternatives engine (manual + fallback logic)
-   Similarity engine
-   Graph data API
-   Search API

------------------------------------------------------------------------

### Frontend Features

#### 1. Tool Page

-   Tool details
-   Alternatives (with explanations)
-   Similar tools
-   Relationships
-   Mini graph

#### 2. Alternatives Page

-   List of alternatives
-   Explanation for each

#### 3. Graph Page

-   Interactive visualization of ecosystem

#### 4. Search

-   Autocomplete
-   Direct navigation

------------------------------------------------------------------------

## 🧱 Data Model

### Tool

type Tool = { id: string name: string slug: string category: "language"
\| "framework" \| "library" subType?: string language?: string
releaseYear?: number status?: string description?: string useCases:
string\[\] tags?: string\[\] website?: string github?: string createdAt:
string updatedAt: string }

------------------------------------------------------------------------

### Relationship

type Relationship = { id: string from: string to: string type: string
weight?: number metadata?: { reason?: string } createdAt: string }

------------------------------------------------------------------------

## 🌐 API Routes

### Tools

-   GET /tools/:id
-   GET /tools/:id/alternatives
-   GET /tools/:id/similar
-   GET /tools/:id/relationships

### Search

-   GET /search?q=react

### Graph

-   GET /tools/:id/graph

------------------------------------------------------------------------

## 🖥️ Frontend Routes

-   / → homepage (search-first)
-   /tools/:slug → tool page
-   /alternatives/:slug → alternatives page
-   /graph/:slug → graph visualization
-   /search?q= → search results

------------------------------------------------------------------------

## 🧠 Alternatives Logic

1.  Manual relationships (ALTERNATIVE_TO)
2.  Fallback:
    -   same subtype
    -   same language
    -   shared useCases
3.  Ranking

------------------------------------------------------------------------

## 🛠️ Tech Stack

### Backend

-   Golang
-   Chi
-   SQLite (better-sqlite3)
-   Goose

### Frontend

-   TanStack Start
-   Tailwind CSS

### Visualization

-   react-force-graph

### Analytics

-   Plausible or Umami

------------------------------------------------------------------------

## ⚙️ Project Setup

### Step 1 --- Initialize project

-   Setup monorepo or separate frontend/backend
-   Install dependencies

### Step 2 --- Setup database

-   Create SQLite DB
-   Create tables
-   Add indexes

### Step 3 --- Seed data

-   Add languages
-   Add tools
-   Add relationships

------------------------------------------------------------------------

## 🔧 Backend Implementation

-   Create DB layer
-   Implement Tool queries
-   Implement Relationships queries
-   Implement Alternatives logic
-   Implement Similar logic
-   Build API endpoints

------------------------------------------------------------------------

## 🎨 Frontend Implementation

### Step 1

-   Setup Next.js project
-   Setup Tailwind

### Step 2

-   Build Tool page

### Step 3

-   Build Alternatives section

### Step 4

-   Build Search

### Step 5

-   Build Graph visualization

------------------------------------------------------------------------

## 🧪 Testing

-   Test API endpoints
-   Validate alternatives logic
-   Test UI navigation

------------------------------------------------------------------------

## 🚀 Deployment

-   Deploy backend (Fly.io / Railway)
-   Deploy frontend (Vercel)

------------------------------------------------------------------------

## 📈 SEO Setup

-   Generate pages:
    -   Alternatives to React
    -   Tools for Python
-   Add meta tags
-   Add internal linking

------------------------------------------------------------------------

## 📊 Analytics

-   Track:
    -   page views
    -   search queries
    -   popular tools

------------------------------------------------------------------------

## ✅ To-Do Checklist

### Data

-   [ ] Define tools dataset
-   [ ] Define relationships
-   [ ] Add explanations

### Backend

-   [ ] Setup SQLite
-   [ ] Create schema
-   [ ] Seed database
-   [ ] Implement API
-   [ ] Implement alternatives logic

### Frontend

-   [ ] Setup Next.js
-   [ ] Build Tool page
-   [ ] Build Alternatives UI
-   [ ] Add Search
-   [ ] Add Graph

### UX

-   [ ] Add explanations
-   [ ] Improve navigation
-   [ ] Handle empty states

### Launch

-   [ ] Create SEO pages
-   [ ] Deploy app
-   [ ] Share on Hacker News
-   [ ] Share on Indie Hackers

------------------------------------------------------------------------

## 💥 Final Strategy

> Ship fast, focus on alternatives + graph, iterate with real user
> feedback.
