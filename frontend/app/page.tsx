"use client";

import ItemsTable, { Item } from "./components/ItemsTable";

// Mock data - replace this with actual database fetching
const mockItems: Item[] = [
  { id: 1, name: "Pepsi Max", price: 2.49 },
  { id: 2, name: "Vitamin Well Refresh", price: 3.29 },
  { id: 3, name: "Estrella Maapähkinä Rinkula", price: 2.99 },
  { id: 4, name: "Red Bull", price: 2.95 },
];

// Current active basket ID
const currentBasketId = 1;

export default function Home() {
  const handleCreateBasket = () => {
    // This will be connected to your API
    console.log("Creating new basket...");
    alert("Create basket functionality - connect to your API");
  };

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
      {/* Header */}
      <header className="bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-zinc-900 dark:text-zinc-50">Walk Through™</h1>
              <p className="text-zinc-600 dark:text-zinc-400 mt-1">The future of shopping</p>
            </div>

            {/* Current Basket Display */}
            <div className="flex items-center gap-4">
              <div className="bg-zinc-100 dark:bg-zinc-800 rounded-lg px-4 py-3">
                <p className="text-sm text-zinc-600 dark:text-zinc-400">Current Basket ID</p>
                <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">{currentBasketId}</p>
              </div>
              <button
                onClick={handleCreateBasket}
                className="px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
              >
                Create New Basket
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="space-y-8">
          {/* Summary Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
              <h3 className="text-sm font-medium text-zinc-600 dark:text-zinc-400">Total Items</h3>
              <p className="text-3xl font-bold text-zinc-900 dark:text-zinc-50 mt-2">{mockItems.length}</p>
              <p className="text-sm text-zinc-500 dark:text-zinc-500 mt-1">Available products</p>
            </div>

            <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
              <h3 className="text-sm font-medium text-zinc-600 dark:text-zinc-400">Total Value</h3>
              <p className="text-3xl font-bold text-zinc-900 dark:text-zinc-50 mt-2">
                €{mockItems.reduce((sum, item) => sum + item.price, 0).toFixed(2)}
              </p>
              <p className="text-sm text-zinc-500 dark:text-zinc-500 mt-1">Combined item prices</p>
            </div>
          </div>

          {/* Items Table */}
          <ItemsTable items={mockItems} />
        </div>
      </main>
    </div>
  );
}
