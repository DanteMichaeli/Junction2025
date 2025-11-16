"use client";

import { useState, useRef } from "react";

interface ClassificationResult {
  itemId: string;
  itemName: string;
  price: number;
  confidence: number;
  matched: boolean;
}

interface ImageClassifierProps {
  onItemClassified?: (result: ClassificationResult) => void;
}

const API_URL = "https://walkthrough-backend-719447017050.europe-north1.run.app";

export default function ImageClassifier({ onItemClassified }: ImageClassifierProps) {
  const [isClassifying, setIsClassifying] = useState(false);
  const [result, setResult] = useState<ClassificationResult | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Create preview
    const reader = new FileReader();
    reader.onload = (e) => {
      setPreviewUrl(e.target?.result as string);
    };
    reader.readAsDataURL(file);

    // Classify the image
    await classifyImage(file);
  };

  const classifyImage = async (file: File) => {
    setIsClassifying(true);
    setError(null);
    setResult(null);

    try {
      const arrayBuffer = await file.arrayBuffer();
      const blob = new Blob([arrayBuffer]);

      const response = await fetch(`${API_URL}/classify-item`, {
        method: "POST",
        body: blob,
        headers: {
          "Content-Type": "application/octet-stream",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to classify image");
      }

      const classificationResult: ClassificationResult = await response.json();
      setResult(classificationResult);

      if (onItemClassified && classificationResult.matched) {
        onItemClassified(classificationResult);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to classify image");
    } finally {
      setIsClassifying(false);
    }
  };

  return (
    <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">Scan Item</h2>

      <div className="space-y-4">
        {/* Preview - Shows after image is uploaded */}
        {previewUrl && (
          <div className="relative">
            <img
              src={previewUrl}
              alt="Preview"
              className="w-full rounded-lg max-h-96 object-contain bg-zinc-100 dark:bg-zinc-800"
            />
          </div>
        )}

        {/* Upload Button */}
        <button
          type="button"
          onClick={() => fileInputRef.current?.click()}
          disabled={isClassifying}
          className="w-full px-8 py-6 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white font-bold text-xl rounded-lg transition-colors shadow-lg hover:shadow-xl flex items-center justify-center gap-3"
        >
          <span className="text-3xl">📁</span>
          {isClassifying ? "Classifying..." : "UPLOAD & CLASSIFY IMAGE"}
        </button>

        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileSelect}
          className="hidden"
        />

        {/* Classification Result */}
        {result && (
          <div
            className={`p-4 rounded-lg ${
              result.matched
                ? "bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800"
                : "bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800"
            }`}
          >
            {result.matched ? (
              <div>
                <h3 className="font-semibold text-green-900 dark:text-green-100 mb-2">Item Identified!</h3>
                <p className="text-green-800 dark:text-green-200">
                  <strong>{result.itemName}</strong>
                </p>
                <p className="text-green-700 dark:text-green-300 text-sm">
                  Price: €{result.price.toFixed(2)}
                </p>
                <p className="text-green-600 dark:text-green-400 text-xs mt-1">
                  Confidence: {(result.confidence * 100).toFixed(1)}%
                </p>
              </div>
            ) : (
              <div>
                <h3 className="font-semibold text-yellow-900 dark:text-yellow-100 mb-2">
                  Item Not Recognized
                </h3>
                <p className="text-yellow-800 dark:text-yellow-200 text-sm">
                  Please try again with a clearer image or a different angle.
                </p>
              </div>
            )}
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="p-4 rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800">
            <p className="text-red-800 dark:text-red-200 text-sm">{error}</p>
          </div>
        )}
      </div>
    </div>
  );
}
