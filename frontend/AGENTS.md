<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.

# Project Guideline:

This project uses Next.js, React, TypeScript, SCSS Modules, TanStack Query, and TanStack Table.

## Core Principles

Prefer simple, explicit code over clever abstractions.

Separate:
- UI rendering
- server/client data fetching
- table behavior
- local UI state
- styling

Use feature-based structure. Domain code belongs in `features`. Reusable UI belongs in `components/ui`.

## Next.js Rendering

Use Server Components for:
- SEO-critical public pages
- metadata
- auth checks
- redirects
- reading cookies
- static or mostly read-only route content

Use Client Components for:
- interactivity
- forms
- browser APIs
- TanStack Query
- TanStack Table
- local UI state

Do not add `"use client"` everywhere by default. Keep the client boundary as small as practical.

## Data Fetching

Use server fetching when initial HTML matters for SEO.

Use TanStack Query for:
- authenticated app screens
- interactive data
- filters
- search
- pagination
- sorting
- mutations
- cache invalidation
- background refetching

For public interactive pages, fetch initial SEO content on the server, then use client components for interactivity.

## TanStack Query

API files perform requests.
Query hooks bind API functions to TanStack Query.
Components consume query hooks.

Do not call `fetch` or `axios` directly inside components unless there is a strong reason.

Use query key factories per feature.
Include every query variable in the query key.
Every mutation should invalidate or update affected queries.

Every query-driven screen must handle:
- loading
- error
- empty
- success

## TanStack Table

TanStack Table manages table behavior, not data fetching or styling.

Keep fetching outside generic table components.
Define columns in separate files.
Use accessor columns for real data fields.
Use display columns for actions and UI-only cells.
Display columns must have explicit IDs.

Put shareable table state in the URL:
- page
- search
- filters
- sorting

Keep non-shareable state local:
- row selection
- column visibility
- open menus

Use pagination or virtualization for large tables.

## State

Use:
- TanStack Query for server/cache state
- TanStack Table for table state
- URL params for shareable filters/search/page/sort
- `useState` for temporary UI state
- `useReducer` for complex local transitions

Avoid global state unless multiple distant features truly need the same client-only state.

Do not store derived values in state.

## TypeScript

Avoid `any`.
Use `unknown` when the type is truly unknown.

Use `interface` for object shapes and component props.
Use `type` for unions, primitives, mapped types, and utility compositions.

Keep feature-specific types close to the feature.

## Components

Page components compose.
Feature components know the domain.
Shared UI components should not know the domain.

Avoid deeply nested JSX.
Use early returns or extracted components.

Use `onSomething` for component props.
Use `handleSomething` for internal handlers.

## Styling

Use SCSS Modules for component styles.
Use global SCSS only for reset, tokens, typography defaults, and base styles.

Use CSS variables for repeated design tokens:
- colors
- spacing
- radius
- shadows
- typography
- z-index

Limit SCSS nesting to 2 levels where possible.
Avoid long parent-child selector chains.

## Performance

Keep Client Components as small as practical.
Avoid unnecessary `useEffect`.
Avoid unnecessary state.

Do not add `useMemo` or `useCallback` everywhere by default.
Use them only for real performance or referential-stability needs.

Keep column definitions outside render.
Avoid expensive table cell renderers.
Use dynamic imports for heavy client-only components.
Use Next.js Image where applicable.

## Naming

Components: `PascalCase`
Hooks: `useSomething`
Files/folders: `kebab-case`
SCSS modules: `component-name.module.scss`

API functions:
- `getUsers`
- `getUserById`
- `createUser`
- `updateUser`
- `deleteUser`

Query hooks:
- `useUsersQuery`
- `useUserQuery`
- `useCreateUserMutation`
- `useUpdateUserMutation`
<!-- END:nextjs-agent-rules -->
