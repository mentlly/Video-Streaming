'use client';

import { useState } from 'react';

export default function VideoUpload() {
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('');

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return alert('Please select a video file first.');
    if (!title.trim()) return alert('Please enter a title.');

    const formData = new FormData();
    
    // Append the file and user-inputted text data
    formData.append('video', file);
    formData.append('title', title);
    formData.append('description', description);

    try {
      setStatus('Uploading...');

      const response = await fetch('/api/video/upload', {
        method: 'POST',
        body: formData,
      });

      if (response.status === 200) {
        setStatus('Upload successful! 🎉');
        // Optional: Clear form on success
        setTitle('');
        setDescription('');
        setFile(null);
      } else {
        const errorText = await response.text();
        setStatus(errorText || 'Upload failed.');
      }
    } catch (error) {
      console.error(error);
      setStatus('An error occurred.');
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 p-4 dark:bg-gray-900">
      <div className="w-full max-w-md rounded-xl bg-white p-6 shadow-md dark:bg-gray-800">
        <h2 className="mb-6 text-2xl font-bold text-gray-900 dark:text-white">
          Upload Video
        </h2>
        
        <form onSubmit={handleUpload} className="space-y-5">
          {/* Title Input */}
          <div>
            <label className="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
              Video Title
            </label>
            <input
              type="text"
              placeholder="e.g., My Awesome Video"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-950 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-500"
              required
            />
          </div>

          {/* Description Input */}
          <div>
            <label className="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
              Description
            </label>
            <textarea
              placeholder="Tell viewers about your video..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-950 placeholder-gray-400 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-700 dark:text-white dark:placeholder-gray-500"
            />
          </div>

          {/* File input */}
          <div>
            <label className="mb-1.5 block text-sm font-medium text-gray-700 dark:text-gray-300">
              Video File
            </label>
            <input 
              type="file" 
              accept="video/*" 
              onChange={(e) => setFile(e.target.files?.[0] || null)} 
              className="w-full text-sm text-gray-500 file:mr-4 file:rounded-lg file:border-0 file:bg-blue-50 file:px-4 file:py-2 file:text-sm file:font-semibold file:text-blue-700 hover:file:bg-blue-100 dark:text-gray-400 dark:file:bg-gray-700 dark:file:text-blue-400"
            />
          </div>

          {/* Submit Button */}
          <button 
            type="submit"
            className="w-full rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:bg-blue-500 dark:hover:bg-blue-600"
          >
            Upload Video
          </button>

          {/* Status Message */}
          {status && (
            <p className={`text-center text-sm font-medium ${
              status.includes('successful') ? 'text-green-600 dark:text-green-400' : 
              status.includes('Uploading') ? 'text-blue-600 dark:text-blue-400' : 'text-red-600 dark:text-red-400'
            }`}>
              {status}
            </p>
          )}
        </form>
      </div>
    </div>
  );
}