"use client";

import { useEffect, useState } from "react";
import ItemsTable, { Item } from "./components/ItemsTable";

// Backend API URL
const API_URL = "http://localhost:3001";

export default function Home() {
  const [items, setItems] = useState<Item[]>([]); // Store catalog
  const [cartItems, setCartItems] = useState<Item[]>([]); // Items in cart
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

  // Set up SSE for real-time updates (items added to cart)
  useEffect(() => {
    const eventSource = new EventSource(`${API_URL}/events`);

    eventSource.onmessage = (event) => {
      try {
        const newItem: Item = JSON.parse(event.data);
        console.log("New item added to cart:", newItem);
        
        // Add item to cart (allow duplicates)
        setCartItems((prevCart) => [...prevCart, newItem]);
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

  // Test function to add all items to basket
  const handleTestAddItems = async () => {
    const itemIds = ["pepsi-max", "sunmaid-sour-raisins", "vitamin-well-refresh", "estrella-chips"];
    
    for (const itemId of itemIds) {
      try {
        const response = await fetch(`${API_URL}/add-item-to-basket`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ itemId }),
        });
        
        if (response.ok) {
          console.log(`✅ Added ${itemId} to basket`);
        }
      } catch (error) {
        console.error(`❌ Failed to add ${itemId}:`, error);
      }
      
      // Small delay between requests
      await new Promise(resolve => setTimeout(resolve, 200));
    }
  };

  // Reset demo function
  const handleResetDemo = async () => {
    try {
      const response = await fetch(`${API_URL}/reset-demo`, {
        method: "POST",
      });
      
      if (response.ok) {
        console.log("✅ Demo reset successfully");
        // Clear cart items in frontend
        setCartItems([]);
        alert("Demo reset! Cart cleared.");
      }
    } catch (error) {
      console.error("❌ Failed to reset demo:", error);
      alert("Failed to reset demo");
    }
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
            
            {/* Test Buttons */}
            <div className="flex gap-2">
              <button
                onClick={handleTestAddItems}
                className="px-4 py-2 bg-purple-600 hover:bg-purple-700 text-white font-medium rounded-lg transition-colors text-sm"
              >
                🧪 Test: Add All Items
              </button>
              <button
                onClick={handleResetDemo}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors text-sm"
              >
                🔄 Reset Demo
              </button>
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
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Left Side: Shopping Cart */}
            <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                  Cart
                </h2>
                <span className="text-sm text-zinc-600 dark:text-zinc-400">
                  Cart ID: {currentBasketId}
                </span>
              </div>
              
              {cartItems.length === 0 ? (
                <div className="text-center py-12 text-zinc-500 dark:text-zinc-400">
                  Cart is empty. Scan groceries with AR spectacles to add them!
                </div>
              ) : (
                <div className="space-y-3">
                  {cartItems.map((item, index) => (
                    <div
                      key={`${item.id}-${index}`}
                      className="flex items-center justify-between p-3 bg-zinc-50 dark:bg-zinc-800 rounded-lg"
                    >
                      <div>
                        <p className="font-medium text-zinc-900 dark:text-zinc-50">{item.name}</p>
                        <p className="text-xs text-zinc-500 dark:text-zinc-400">#{item.id}</p>
                      </div>
                      <p className="font-bold text-zinc-900 dark:text-zinc-50">
                        €{item.price.toFixed(2)}
                      </p>
                    </div>
                  ))}
                  
                  {/* Cart Total */}
                  <div className="pt-3 mt-3 border-t border-zinc-200 dark:border-zinc-700">
                    <div className="flex items-center justify-between">
                      <p className="text-lg font-bold text-zinc-900 dark:text-zinc-50">Total:</p>
                      <p className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                        €{cartItems.reduce((sum, item) => sum + item.price, 0).toFixed(2)}
                      </p>
                    </div>
                    <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                      {cartItems.length} item{cartItems.length !== 1 ? 's' : ''} in cart
                    </p>
                  </div>
                </div>
              )}
            </div>

            {/* Right Side: Store Items */}
            <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
              <h2 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50 mb-4">
                Store Items
              </h2>
              <ItemsTable items={items} />
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
