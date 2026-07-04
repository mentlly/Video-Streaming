import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  reactCompiler: true,
  logging: {
    fetches: {
      fullUrl: true,
    },
  },
  allowedDevOrigins: ['192.168.1.*'],
};

export default nextConfig;
