"use server"

const BACKEND_URL = process.env.BACKEND_URL;

export interface videoDetails {
    title: string,
    videoId: string,
    description: string,
    duration: number,
    thumbnailUrl: string,
    channelName: string,
}

export async function getLatestVideos(): Promise<videoDetails[]> {
    const response = await fetch(BACKEND_URL+"/api/video/get/latest", {
        method: "GET",
        headers: { 'Content-Type': 'application/json' }
    });
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    const data: videoDetails[] = await response.json();
    console.log(data);
    return data;
}