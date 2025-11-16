"use client";

import { useEffect, useState } from "react";

interface LeaderboardEntry {
  ownerName: string;
  durationSecs: number;
  completedAt: string;
}

interface LeaderboardProps {
  apiUrl: string;
}

export default function Leaderboard({ apiUrl }: LeaderboardProps) {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchLeaderboard = async () => {
      try {
        const response = await fetch(`${apiUrl}/leaderboard`);
        if (response.ok) {
          const data = await response.json();
          setEntries(data || []);
        }
      } catch (error) {
        console.error("Error fetching leaderboard:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchLeaderboard();

    // Refresh leaderboard every 5 seconds
    const interval = setInterval(fetchLeaderboard, 5000);
    return () => clearInterval(interval);
  }, [apiUrl]);

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50 mb-4">
        🏆 Leaderboard
      </h2>
      
      {isLoading ? (
        <div className="text-center py-8 text-zinc-500 dark:text-zinc-400">
          Loading...
        </div>
      ) : entries.length === 0 ? (
        <div className="text-center py-8 text-zinc-500 dark:text-zinc-400">
          No completed baskets yet. Be the first!
        </div>
      ) : (
        <div className="space-y-2 max-h-60 overflow-y-auto pr-2">
          {entries.map((entry, index) => (
            <div
              key={index}
              className={`flex items-center justify-between p-3 rounded-lg ${
                index === 0
                  ? 'bg-yellow-50 dark:bg-yellow-900/20 border-2 border-yellow-400'
                  : index === 1
                  ? 'bg-zinc-100 dark:bg-zinc-800 border-2 border-zinc-400'
                  : index === 2
                  ? 'bg-orange-50 dark:bg-orange-900/20 border-2 border-orange-400'
                  : 'bg-zinc-50 dark:bg-zinc-800'
              }`}
            >
              <div className="flex items-center gap-3">
                <span className="text-2xl font-bold text-zinc-600 dark:text-zinc-400">
                  {index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : `#${index + 1}`}
                </span>
                <p className="font-semibold text-zinc-900 dark:text-zinc-50">
                  {entry.ownerName}
                </p>
              </div>
              <p className="text-3xl font-mono font-bold text-zinc-900 dark:text-zinc-50">
                {entry.durationSecs}s
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

