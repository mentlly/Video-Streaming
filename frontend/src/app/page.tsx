'use client';

import { useState } from 'react';

// Mock data to simulate videos coming from your database/API
const MOCK_VIDEOS = [
  {
    id: '1',
    title: 'Building a Next.js App in 2026: The Ultimate Guide',
    channelName: 'CodeCraft',
    thumbnailUrl: 'https://images.unsplash.com/photo-1618401471353-b98afee0b2eb?w=500&auto=format&fit=crop&q=60',
    duration: '14:22',
  },
  {
    id: '2',
    title: 'Learn Tailwind CSS in 10 Minutes flat',
    channelName: 'DesignBytes',
    thumbnailUrl: 'https://images.unsplash.com/photo-1507238691740-187a5b1d37b8?w=500&auto=format&fit=crop&q=60',
    duration: '10:00',
  },
  {
    id: '3',
    title: 'Why Everyone is Switching to Full-Stack Development',
    channelName: 'TechPulse',
    thumbnailUrl: 'https://images.unsplash.com/photo-1531403009284-440f080d1e12?w=500&auto=format&fit=crop&q=60',
    duration: '22:15',
  },
];

export default function VideoFeed() {
  const [videos] = useState(MOCK_VIDEOS);

  const handleVideoClick = (videoId: string) => {
    console.log(`Navigating to video / watch page for ID: ${videoId}`);
    // If using Next.js router: router.push(`/watch/${videoId}`)
  };

  return (
    <div className="min-h-screen bg-gray-50 px-4 py-8 dark:bg-gray-900 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <h1 className="mb-8 text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
          Recommended Videos
        </h1>

        {/* Responsive Grid Layout */}
        <div className="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {videos.map((video) => (
            <div
              key={video.id}
              onClick={() => handleVideoClick(video.id)}
              className="group cursor-pointer overflow-hidden rounded-xl bg-white transition-all duration-200 hover:shadow-md dark:bg-gray-800"
            >
              {/* Thumbnail Container */}
              <div className="relative aspect-video w-full overflow-hidden bg-gray-200 dark:bg-gray-700">
                <img
                  src={video.thumbnailUrl}
                  alt={video.title}
                  className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
                />
                {/* Optional Duration Tag */}
                <span className="absolute bottom-2 right-2 rounded bg-black/75 px-1.5 py-0.5 text-xs font-medium text-white">
                  {video.duration}
                </span>
              </div>

              {/* Video Info Details */}
              <div className="p-4">
                <h3 className="line-clamp-2 text-sm font-semibold leading-tight text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400">
                  {video.title}
                </h3>
                <p className="mt-1.5 text-xs font-medium text-gray-500 dark:text-gray-400">
                  {video.channelName}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}