"use client";

import { useEffect, useState } from "react";
import ItemsTable, { Item } from "./components/ItemsTable";
import Leaderboard from "./components/Leaderboard";

// Backend API URL
const API_URL = "http://localhost:3001";

export default function Home() {
  const [items, setItems] = useState<Item[]>([]); // Store catalog
  const [cartItems, setCartItems] = useState<Item[]>([]); // Items in cart
  const [currentBasketId, setCurrentBasketId] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [userName, setUserName] = useState<string>("");
  const [hasStarted, setHasStarted] = useState(false);
  const [nameInput, setNameInput] = useState<string>("");
  const [startTime, setStartTime] = useState<number>(0);
  const [elapsedTime, setElapsedTime] = useState<number>(0);
  const [isComplete, setIsComplete] = useState(false);

  // Fetch items and basket ID from backend
  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch items
        const itemsResponse = await fetch(`${API_URL}/items`);
        if (itemsResponse.ok) {
          const data = await itemsResponse.json();
          setItems(data || []);
        }

        // Fetch current basket ID
        const basketResponse = await fetch(`${API_URL}/current-basket`);
        if (basketResponse.ok) {
          const basketData = await basketResponse.json();
          setCurrentBasketId(basketData.basketId);
        }
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  // Timer effect - updates every second
  useEffect(() => {
    if (!hasStarted || isComplete) return;

    const interval = setInterval(() => {
      setElapsedTime(Date.now() - startTime);
    }, 1000);

    return () => clearInterval(interval);
  }, [hasStarted, startTime, isComplete]);

  // Set up SSE for real-time updates (items added to cart)
  useEffect(() => {
    const eventSource = new EventSource(`${API_URL}/events`);

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        
        // Check if this is an item or completion message
        if (data.item) {
          const newItem: Item = data.item;
          console.log("New item added to cart:", newItem);
          
          // Add item to cart (allow duplicates)
          setCartItems((prevCart) => [...prevCart, newItem]);
          
          // Check if basket is complete
          if (data.isComplete) {
            setIsComplete(true);
            console.log("🎉 Basket completed!");
          }
        } else {
          // Legacy format (just item)
          const newItem: Item = data;
          setCartItems((prevCart) => [...prevCart, newItem]);
        }
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
    const itemIds = ["red-bull", "vitamin-well-refresh", "estrella-chips"];
    
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
        
        // Clear everything and return to welcome screen
        setCartItems([]);
        setCurrentBasketId("");
        setHasStarted(false);
        setUserName("");
        setNameInput("");
        
        alert("Demo reset! Returning to welcome screen.");
      }
    } catch (error) {
      console.error("❌ Failed to reset demo:", error);
      alert("Failed to reset demo");
    }
  };

  // Format elapsed time as MM:SS
  const formatTime = (ms: number) => {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
  };

  // Handle start shopping
  const handleStartShopping = async () => {
    if (!nameInput.trim()) return;
    
    try {
      // Create new basket with user's name
      const response = await fetch(`${API_URL}/create-basket`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ ownerName: nameInput.trim() }),
      });
      
      if (response.ok) {
        const data = await response.json();
        console.log("✅ Basket created:", data);
        
        // Set user name and basket ID
        setUserName(nameInput.trim());
        setCurrentBasketId(data.basketId);
        
        // Start the timer
        setStartTime(Date.now());
        setElapsedTime(0);
        setIsComplete(false);
        
        setHasStarted(true);
      } else {
        alert("Failed to create basket");
      }
    } catch (error) {
      console.error("❌ Failed to create basket:", error);
      alert("Failed to create basket");
    }
  };

  // Show welcome screen if user hasn't started
  if (!hasStarted) {
    return (
      <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center">
        <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-xl p-8 max-w-md w-full">
          <h1 className="text-4xl font-bold text-zinc-900 dark:text-zinc-50 mb-2 text-center">
            Walk Through™
          </h1>
          <p className="text-zinc-600 dark:text-zinc-400 mb-8 text-center">
            The future of shopping
          </p>
          
          <div className="space-y-4">
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                Enter your name to start shopping
              </label>
              <input
                id="name"
                type="text"
                value={nameInput}
                onChange={(e) => setNameInput(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleStartShopping()}
                placeholder="Your name"
                className="w-full px-4 py-3 border border-zinc-300 dark:border-zinc-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
              />
            </div>
            
            <button
              onClick={handleStartShopping}
              disabled={!nameInput.trim()}
              className="w-full px-6 py-4 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-bold text-lg rounded-lg transition-colors shadow-lg hover:shadow-xl"
            >
              Start Shopping 🛒
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
      {/* Header */}
      <header className="bg-white dark:bg-zinc-900 border-b border-zinc-200 dark:border-zinc-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-zinc-900 dark:text-zinc-50">Walk Through™</h1>
              <p className="text-zinc-600 dark:text-zinc-400 mt-1">Welcome, {userName}!</p>
            </div>
            
            {/* Timer Display */}
            <div className="flex items-center gap-4">
              <div className={`px-6 py-3 rounded-lg ${
                isComplete 
                  ? 'bg-green-100 dark:bg-green-900 border-2 border-green-500' 
                  : 'bg-blue-100 dark:bg-blue-900'
              }`}>
                <p className="text-xs text-zinc-600 dark:text-zinc-400 mb-1">
                  {isComplete ? '✅ Completed' : '⏱️ Time'}
                </p>
                <p className={`text-3xl font-mono font-bold ${
                  isComplete 
                    ? 'text-green-700 dark:text-green-300' 
                    : 'text-blue-700 dark:text-blue-300'
                }`}>
                  {formatTime(elapsedTime)}
                </p>
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

        {/* Leaderboard Section */}
        {!isLoading && (
          <div className="mt-8">
            <Leaderboard apiUrl={API_URL} />
          </div>
        )}
      </main>
    </div>
  );
}
