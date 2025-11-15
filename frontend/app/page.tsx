"use client";

import { useEffect, useState } from "react";
import ItemsTable, { Item } from "./components/ItemsTable";
import ImageClassifier from "./components/ImageClassifier";

// Backend API URL
const API_URL = "http://localhost:3001";

export default function Home() {
  const [items, setItems] = useState<Item[]>([]);
  const [currentBasketId] = useState(1);
  const [isLoading, setIsLoading] = useState(true);

  // Fetch items from backend
  useEffect(() => {
    const fetchItems = async () => {
      try {
        const response = await fetch(`${API_URL}/items`);
        if (response.ok) {
          const data = await response.json();
          setItems(data || []);
        }
      } catch (error) {
        console.error("Error fetching items:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchItems();
  }, []);

  // Set up SSE for real-time updates
  useEffect(() => {
    const eventSource = new EventSource(`${API_URL}/events`);

    eventSource.onmessage = (event) => {
      try {
        const newItem: Item = JSON.parse(event.data);
        console.log("New item created:", newItem);
        
        // Add new item to the list if it doesn't already exist
        setItems((prevItems) => {
          const exists = prevItems.some((item) => item.id === newItem.id);
          if (!exists) {
            return [...prevItems, newItem];
          }
          return prevItems;
        });
      } catch (error) {
        console.error("Error parsing SSE data:", error);
      }
    };

    eventSource.onerror = (error) => {
      console.error("SSE connection error:", error);
    };

    // Cleanup on unmount
    return () => {
      eventSource.close();
    };
  }, []);

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
            <div className="bg-zinc-100 dark:bg-zinc-800 rounded-lg px-4 py-3">
              <p className="text-sm text-zinc-600 dark:text-zinc-400">Current Basket ID</p>
              <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">{currentBasketId}</p>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {isLoading ? (
          <div className="flex items-center justify-center py-12">
            <div className="text-zinc-600 dark:text-zinc-400">Loading items...</div>
          </div>
        ) : (
          <div className="space-y-8">
            {/* Image Classifier */}
            <ImageClassifier
              onItemClassified={(result) => {
                console.log("Item classified:", result);
                // Optionally add to basket or show notification
              }}
            />

            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
                <h3 className="text-sm font-medium text-zinc-600 dark:text-zinc-400">Total Items</h3>
                <p className="text-3xl font-bold text-zinc-900 dark:text-zinc-50 mt-2">{items.length}</p>
                <p className="text-sm text-zinc-500 dark:text-zinc-500 mt-1">Available products</p>
              </div>

              <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
                <h3 className="text-sm font-medium text-zinc-600 dark:text-zinc-400">Total Value</h3>
                <p className="text-3xl font-bold text-zinc-900 dark:text-zinc-50 mt-2">
                  €{items.reduce((sum, item) => sum + item.price, 0).toFixed(2)}
                </p>
                <p className="text-sm text-zinc-500 dark:text-zinc-500 mt-1">Combined item prices</p>
              </div>
            </div>

            {/* Items Table */}
            <ItemsTable items={items} />
          </div>
        )}
      </main>
    </div>
  );
}
