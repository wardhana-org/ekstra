import { auth } from "~/server/auth";

export default auth((req) => {
  const { pathname } = req.nextUrl;
  const creatorPageId = req.auth?.user?.creatorPageId ?? null; // Explicit null coalescing

  // Route definitions
  const publicRoutes = ['/', '/login', '/about', '/gettingstarted'];
  const creatorRoutes = ['/creator', '/creator/*']; // Include all sub-routes
  
  // Route checks
  const isPublicRoute = publicRoutes.includes(pathname);
  const isCreatorRoute = creatorRoutes.some(route => 
    pathname.startsWith(route.replace('/*', ''))
  );

  // 1. First handle creator routes - highest priority
  if (isCreatorRoute) {
    if (!req.auth) {
      // Not logged in at all - send to login
      return Response.redirect(new URL("/login", req.nextUrl.origin));
    }
    if (creatorPageId === null) {
      // Logged in but not a creator - send to home
      return Response.redirect(new URL("/user/home", req.nextUrl.origin));
    }
    // Valid creator (creatorPageId exists) - allow access
    return;
  }

  // 2. Handle public routes
  if (isPublicRoute) {
    if (req.auth) {
      // Logged in but on public route - send to home
      return Response.redirect(new URL("/user/home", req.nextUrl.origin));
    }
    // Not logged in on public route - allow access
    return;
  }

  // 3. All other routes (protected non-creator routes)
  if (!req.auth) {
    // Not logged in - send to login
    return Response.redirect(new URL("/login", req.nextUrl.origin));
  }

  // Logged in user accessing non-public, non-creator route - allow access
  return;
});

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};