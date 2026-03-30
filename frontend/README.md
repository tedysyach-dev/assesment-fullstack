Stack: Vite + React + TypeScript + Tailwind v4 + React Router + Axios + Zustand
Struktur folder feature-based:

```bash
src/
├── features/
│ ├── auth/
│ │ ├── components/ → LoginPage, ProtectedRoute
│ │ ├── hooks/ → useAuth.ts
│ │ ├── store/ → authStore.ts (Zustand + persist)
│ │ └── types/ → index.ts
│ └── home/
│ └── components/ → HomePage.tsx
├── lib/
│ └── axios.ts → instance + interceptors
├── router/
│ └── index.tsx → semua route config
├── components/ui/ → shared UI components
├── hooks/ → shared hooks
└── types/ → shared types
```

Highlights yang sudah tersetup:
Axios dengan request interceptor (auto-attach token) dan response interceptor (handle 401 → redirect login)
Zustand auth store dengan persist middleware — state tetap ada setelah refresh
ProtectedRoute — kalau belum login, auto redirect ke /login
Path alias @/ → src/ di Vite dan tsconfig
