# Setup Instructions

## Backend Setup

### Initialize the Database

1. Navigate to the backend directory:
```bash
cd backend
```

2. Initialize the database (creates tables and adds sample data):
```bash
go run cmd/initdb.go
```

### Start the Backend Server

Run the server on port 3001:
```bash
go run cmd/server.go
```

The backend will be available at `http://localhost:3001`

### Backend Endpoints

- `GET /items` - Get all items
- `GET /items?id={id}` - Get a specific item
- `POST /items` - Create a new item
- `PUT /items` - Update an item
- `DELETE /items?id={id}` - Delete an item
- `GET /events` - SSE endpoint for real-time updates

### Test with cURL

Create an item:
```bash
curl -X POST "http://localhost:3001/items" -H "Content-Type: application/json" -d '{"id":"new-item","name":"New Item","price":12.99}'
```

Get all items:
```bash
curl -X GET "http://localhost:3001/items"
```

Delete an item:
```bash
curl -X DELETE "http://localhost:3001/items?id=new-item"
```

## Frontend Setup

### Install Dependencies

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install npm packages:
```bash
npm install
```

### Start the Frontend Development Server

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## How It Works

1. The frontend fetches all items from the backend on page load
2. The frontend establishes an SSE connection to `/events`
3. When a new item is created via the backend API, all connected frontends receive a real-time update
4. The new item automatically appears in the dashboard without refreshing the page

## Testing Real-Time Updates

1. Start both backend and frontend servers
2. Open the frontend in your browser (`http://localhost:3000`)
3. In a terminal, create a new item using cURL:
```bash
curl -X POST "http://localhost:3001/items" -H "Content-Type: application/json" -d '{"id":"real-time-test","name":"Real-Time Test Item","price":5.99}'
```
4. Watch the item appear in the frontend dashboard automatically!

