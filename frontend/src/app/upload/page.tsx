'use client';

import { useState } from 'react';

export default function VideoUpload() {
  const [file, setFile] = useState<File | null>(null);
  const [status, setStatus] = useState('');

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return alert('Please select a video file first.');

    // 1. Create the FormData instance
    const formData = new FormData();
    
    // 2. Append the file (and any extra text data you want)
    formData.append('video', file);
    formData.append('title', 'My Awesome Video');

    try {
      setStatus('Uploading...');

      const response = await fetch('/api/video/upload', {
        method: 'POST',
        body: formData,
      });

      if (response.status === 200) {
        setStatus('Upload successful!');
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
    <form onSubmit={handleUpload}>
      <input 
        type="file" 
        accept="video/*" 
        onChange={(e) => setFile(e.target.files?.[0] || null)} 
      />
      <button type="submit">Upload Video</button>
      <p>{status}</p>
    </form>
  );
}