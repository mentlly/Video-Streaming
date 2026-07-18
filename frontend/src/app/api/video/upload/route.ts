export const runtime = 'nodejs';

const BACKEND_URL = process.env.BACKEND_URL;

export async function POST(request: Request) {
  const contentType = request.headers.get('content-type') ?? 'multipart/form-data';
  const cookie = request.headers.get('cookie') ?? '';

  const response = await fetch(`${BACKEND_URL}/api/video/upload`, {
    method: 'POST',
    headers: {
      'content-type': contentType,
      ...(cookie ? { 'cookie': cookie } : {}),
    },
    body: request.body,
    duplex: 'half',
  } as RequestInit & { duplex: 'half' });

  const responseText = await response.text();

  return new Response(responseText, {
    status: response.status,
    headers: {
      'content-type': response.headers.get('content-type') ?? 'text/plain; charset=utf-8',
    },
  });
}