import axios from 'axios';

// In production (built via `npm run build`), VITE_API_BASE_URL is set to
// "/api" so requests resolve relative to the CloudFront domain, hitting the
// /api/* cache behavior that forwards to the ALB backend. In local dev it
// falls back to the Go server running directly on :8080.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api';

export const apiClient = axios.create({
    baseURL: API_BASE_URL,
    withCredentials: true, // Crucial for httpOnly cookies
});
