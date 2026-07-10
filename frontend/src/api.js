// IMPORTANT: always call the backend with a RELATIVE path, never a hardcoded
// domain. CloudFront's /api/* behavior forwards these to the ALB origin,
// so the same https://<cloudfront-domain> serves both the SPA and the API
// with no CORS or mixed-content issues.

export async function checkHealth() {
  const res = await fetch("/api/v1/health");
  if (!res.ok) {
    throw new Error(`Health check failed: ${res.status}`);
  }
  return res.json();
}
