import { auth } from "~/server/auth";

export default auth((req) => {
  const { pathname } = req.nextUrl;
  const isPublicRoute = [
    '/',
    '/login',
    '/about',
    '/gettingstarted'
  ].includes(pathname);

  // If user is authenticated redirect to /home
  if (req.auth && isPublicRoute) {
    return Response.redirect(new URL("/user/home", req.nextUrl.origin));
  }

  // If user is not authenticated redirect to /login
  if (!req.auth && !isPublicRoute) {
    return Response.redirect(new URL("/login", req.nextUrl.origin));
  }

  return;
});

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};