# WalkThrough™ - AR Shopping Experience

**The future of shopping** - A real-time AR shopping system using Snapchat Spectacles, built for Junction 2025 hackathon.

## 🎯 Project Overview

WalkThrough is an innovative shopping experience that combines AR glasses (Snapchat Spectacles), real-time backend processing, and a live dashboard to create a seamless, futuristic shopping journey. Users wear AR glasses that automatically identify products, add them to their basket, and track their shopping speed - all displayed on a live leaderboard!

## 🏗️ Architecture

The project consists of three main components:

```
Junction2025/
├── backend/          # Go backend server with SQLite database
├── frontend/         # Next.js dashboard for real-time visualization
└── lens/            # Snapchat Lens Studio AR application
```

### System Flow

```
AR Glasses (Spectacles)
    ↓ (scan and classify product)
    ↓
Backend API (Go + SQLite)
    ↓ (Server-Sent Events)
    ↓
Frontend Dashboard (Next.js)
    ↓ (display in real-time)
    ↓
Leaderboard (fastest shoppers)
```

## 📦 Components

### 1. Backend (Go + SQLite)

**Location:** `/backend`

A RESTful API server that manages:
- Product catalog (3 items: Red Bull, Vitamin Well, Estrella Chips)
- Shopping baskets with owner names and timestamps
- Real-time updates via Server-Sent Events (SSE)
- Automatic basket completion when all 3 items are collected
- Leaderboard tracking

**Tech Stack:**
- Go 1.22+
- SQLite database
- Server-Sent Events for real-time updates
- UUID generation for baskets

**Key Features:**
- Automatic basket completion detection
- Shopping duration tracking (createDate → completedAt)
- Persistent basket history
- Real-time broadcasting to all connected dashboards

### 2. Frontend (Next.js + React)

**Location:** `/frontend`

A real-time dashboard that displays:
- Live shopping cart updates
- Product catalog with images
- Timer showing shopping duration
- Leaderboard of fastest shoppers
- Welcome screen for user registration

**Tech Stack:**
- Next.js 16.0.3
- React with TypeScript
- Tailwind CSS for styling
- Server-Sent Events for real-time updates

**Key Features:**
- Real-time cart updates (no refresh needed)
- Grouped items with quantity badges
- Live timer with completion detection
- Scrollable leaderboard
- Dark mode support

### 3. Lens (Snapchat Spectacles AR)

**Location:** `/lens`

AR application for Snapchat Spectacles that:
- Identifies products using computer vision
- Sends product IDs to backend
- Provides visual feedback to the user

**Tech Stack:**
- Lens Studio
- SnapML for machine learning
- JavaScript for scripting

## 🚀 Quick Start

### Prerequisites

- Go 1.22 or higher
- Node.js 18+ and npm
- Snapchat Lens Studio (for AR development)

### Backend Setup

```bash
cd backend

# Start the server
go run cmd/server.go
```

The backend will:
- Create SQLite database (`app.db`)
- Set up tables (items, baskets, item_basket)
- Insert 3 sample products
- Start server on `http://localhost:3001`

### Frontend Setup

```bash
cd frontend

# Install dependencies (first time only)
npm install

# Start development server
npm run dev
```

The frontend will be available at `http://localhost:3000`

### Lens Setup

1. Open Lens Studio
2. Load the project from `/lens` directory
3. Configure the backend API endpoint
4. Build and deploy to Spectacles

## 📊 Database Schema

### Items Table
```sql
CREATE TABLE items (
    id STRING PRIMARY KEY,
    name STRING NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category STRING NOT NULL,
    thumbnail STRING NOT NULL
);
```

**Sample Data:**
- Red Bull (€2.49) - Beverage
- Vitamin Well Refresh (€2.79) - Beverage
- Estrella Maapähkinä Rinkula (€2.89) - Snacks

### Baskets Table
```sql
CREATE TABLE baskets (
    basketID UUID PRIMARY KEY,
    ownerName TEXT NOT NULL,
    createDate DATETIME NOT NULL,
    completedAt DATETIME
);
```

**Fields:**
- `basketID`: Unique UUID generated for each shopping session
- `ownerName`: Customer's name
- `createDate`: When shopping started (timestamp)
- `completedAt`: When all 3 items were collected (NULL if ongoing)

### Item-Basket Table (Join)
```sql
CREATE TABLE item_basket (
    itemID STRING,
    basketID UUID,
    FOREIGN KEY(itemID) REFERENCES items(id),
    FOREIGN KEY(basketID) REFERENCES baskets(basketID)
);
```

Links items to baskets (many-to-many relationship).
## 🎮 User Flow

### 1. Welcome Screen
- User enters their name
- Clicks "Start Shopping"
- Backend creates new basket with UUID
- Timer starts

### 2. Shopping Experience
- User wears Snapchat Spectacles
- AR glasses identify products
- Items automatically added to basket
- Dashboard updates in real-time
- Timer counts up

### 3. Completion
- When 3 unique items are collected:
  - Timer stops and turns green
  - `completedAt` timestamp recorded
  - "🎉 Basket completed!" logged
  - Duration calculated

### 4. Leaderboard
- Completed baskets appear on leaderboard
- Sorted by duration (fastest first)
- Shows: 🥇🥈🥉 medals for top 3
- Updates every 5 seconds

### 5. Reset
- Click "Reset Demo"
- Returns to welcome screen
- Basket history preserved in database
- Ready for next user

## 🎨 Frontend Features

### Dashboard Layout

**Header:**
- Title and welcome message
- Live timer (MM:SS format)
- Test buttons for demo

**Main Content:**
- **Cart (Left):** Live shopping cart with product images
- **Store Items (Right):** Product catalog with categories
- **Leaderboard (Bottom):** Top 10 fastest shoppers

### Cart Display

Items are grouped by ID with quantity badges:

```
[IMAGE] Red Bull           x2  €4.98
        #red-bull

[IMAGE] Vitamin Well           €2.79
        #vitamin-well-refresh
```

**Features:**
- Product thumbnails (48x48px)
- Quantity badges (x2, x3, etc.)
- Multiplied prices
- Scrollable (max 320px height)
- Running total at bottom

### Leaderboard

```
🥇 Alice Smith      9s
🥈 Bob Johnson     12s
🥉 Charlie Brown   15s
#4 David Lee       18s
```

**Features:**
- Medal icons for top 3
- Smart time formatting (seconds or MM:SS)
- Scrollable (max 240px height)
- Auto-refreshes every 5 seconds

## 🧪 Testing

### Test Buttons

**Test: Add Red Bull**
- Adds one Red Bull to cart
- Click multiple times to test quantity grouping

**Test: Add All Items**
- Adds all 3 items at once
- Triggers basket completion
- Tests full flow

**Reset Demo**
- Clears frontend session
- Returns to welcome screen
- Preserves all database history


## 📈 Data Persistence

### What Persists:
- ✅ All baskets (with owner names)
- ✅ All basket items
- ✅ All completion timestamps
- ✅ Complete shopping history

### What Resets:
- ❌ Active basket ID (cleared on reset)
- ❌ Frontend display (returns to welcome)
- ❌ Timer (resets for new user)

This allows you to:
- Track all demos at the hackathon
- Analyze shopping patterns
- Review basket history
- Calculate statistics

## 🎯 Basket Completion Logic

A basket is considered **complete** when it contains **3 distinct items**.

**Duration Calculation:**
```sql
(julianday(completedAt) - julianday(createDate)) * 86400
```
Result in seconds.

## 🔄 Real-Time Updates (SSE)

### How It Works

1. **Frontend connects** to `/events` endpoint
2. **Backend maintains** list of connected clients
3. **When item added**, backend broadcasts to all clients
4. **Frontend receives** event and updates cart
5. **No page refresh** needed!

### Event Format

```javascript
// Item added
{
  "item": {
    "id": "red-bull",
    "name": "Red Bull",
    "price": 2.49,
    "category": "Beverage",
    "thumbnail": "https://..."
  },
  "isComplete": false
}

// Basket completed
{
  "item": {...},
  "isComplete": true  // ← Triggers timer stop
}
```

## 🛠️ Development

### Project Structure

```
backend/
├── cmd/
│   └── server.go          # Main server entry point
├── pkg/
│   ├── db.go             # Database setup and management
│   └── service.go        # Business logic (items, baskets)
├── go.mod                # Go dependencies
└── app.db               # SQLite database (auto-generated)

frontend/
├── app/
│   ├── components/
│   │   ├── ItemsTable.tsx      # Product catalog display
│   │   └── Leaderboard.tsx     # Leaderboard component
│   ├── page.tsx                # Main dashboard page
│   └── layout.tsx              # App layout
├── package.json
└── next.config.ts

lens/
└── [Lens Studio project files]
```

### Key Files

**Backend:**
- `cmd/server.go` - HTTP server, API endpoints, SSE handling
- `pkg/db.go` - Database schema, table creation, sample data
- `pkg/service.go` - CRUD operations, basket logic, leaderboard

**Frontend:**
- `app/page.tsx` - Main dashboard with cart, timer, controls
- `app/components/ItemsTable.tsx` - Product catalog table
- `app/components/Leaderboard.tsx` - Ranking display

## 📱 AR Glasses Integration

### Lens Studio Setup

The AR glasses (Snapchat Spectacles) run a Lens that:

1. **Captures camera feed**
2. **Runs ML model** to identify products
3. **Sends HTTP request** to backend:
   ```javascript
   POST /add-item-to-basket
   {
     "itemId": "red-bull"
   }
   ```
4. **Receives response** with completion status
5. **Shows visual feedback** to user

### Product Identification

The Lens uses:
- Image classification model
- Trained on 3 product types
- Returns product ID
- Sends to backend API

## 🎯 Demo Workflow

### For Hackathon Judges

1. **Open dashboard** on large screen
2. **Hand Spectacles** to participant
3. **Participant enters name** on dashboard
4. **Timer starts** automatically
5. **Participant scans products** in store
6. **Cart updates in real-time** on screen
7. **Timer stops** when all 3 items found
8. **Leaderboard updates** with new time
9. **Click "Reset Demo"** for next participant

### For Multiple Participants

- Each participant creates their own basket
- All baskets persist in database
- Leaderboard shows all completed baskets
- Competitive element drives engagement!

---

**Built with ❤️ at Junction 2025**

*The future of shopping is here.*

