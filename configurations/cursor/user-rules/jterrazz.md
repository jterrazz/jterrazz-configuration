You are an expert in TypeScript, Node.js, Solidity, Next.js, React, Shadcn UI, Zod, and Tailwind, with a strong commitment to best practices and modern development standards.

Code Style and Structure

- Write clear, concise, and maintainable TypeScript code following SOLID principles and Clean Architecture.
- Prefer functional and declarative programming patterns; use classes sparingly (primarily for explicit use case definitions).
- Promote modularization, iteration, and DRY (Don't Repeat Yourself) principles to avoid duplication.
- Use clear, descriptive variable names with auxiliary verbs (e.g., isLoading, hasError).
- Organize files distinctly into domain interfaces, infrastructure implementations, services, components, helpers, static assets, and type definitions.

Naming Conventions

- Use lowercase with dashes for directory names (e.g., components/auth-wizard).
- Favor named exports over default exports for clearer imports and improved maintainability.
- Avoid abbreviations (e.g., use configuration instead of config).
- Ensure infrastructure implementations are clearly and explicitly named (e.g., configuration.service.ts).

TypeScript Usage

- Use strict TypeScript types consistently across the project.
- Utilize enums thoughtfully for readability and maintainability, especially for domain-specific values.
- Clearly separate domain interfaces (abstract definitions) from infrastructure layers (concrete implementations).

Syntax and Formatting

- Only write code comments useful for the project, not for my prompts.
- Prioritize readability through clear error handling, using guard clauses and early returns.
- Keep conditional logic flat and straightforward by avoiding deeply nested conditions.
- Embrace declarative JSX syntax for React components for clarity and ease of comprehension.

Testing

- Use Jest alongside jest-mock-extended for comprehensive unit and integration testing.
- Clearly structure tests following the Given/When/Then paradigm to enhance clarity and intention.
- Write test descriptions consistently with the "it should" prefix, utilizing root describe blocks for organizational coherence.
- Separate integration tests clearly in the root-level **tests** folder; unit tests reside alongside their respective implementation files.

UI and Styling

- Leverage Shadcn UI and Tailwind CSS exclusively to build consistent and maintainable interfaces.
- Follow responsive, mobile-first design strategies to ensure accessible and intuitive user experiences.

Next.js and React

- Minimize usage of 'use client', 'useEffect', and React state hooks; prefer React Server Components (RSC) to optimize performance.
- Dynamically load client components and wrap them in Suspense with meaningful fallback UI to improve user perception.
- Optimize images proactively using WebP format, including size metadata, and applying lazy loading techniques.

Server Actions and Error Handling

- Implement server actions using next-safe-action with robust Zod schemas for type-safe validation.
- Clearly distinguish expected errors (managed via useActionState) from unexpected errors (handled via dedicated error boundaries like error.tsx and global-error.tsx).
- Ensure service layers consistently throw descriptive, user-friendly errors that can be gracefully managed by tanStackQuery.

Performance and Optimization

- Favor Next.js Server-Side Rendering (SSR) patterns to minimize client-side state complexity.
- Regularly optimize and monitor Web Vitals (LCP, CLS, FID) to ensure performance excellence.

Documentation and Best Practices

- Continuously follow Next.js official documentation and established best practices for data fetching, rendering, routing, and state management.
- Regularly update dependencies and remain aligned with community standards and best practices.
