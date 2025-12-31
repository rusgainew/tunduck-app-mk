/**
 * Dynamic imports for dashboard pages
 * Reduces initial bundle size by lazy loading pages on demand
 */

import { lazy } from 'react';

// Lazy load dashboard pages
export const DocumentsPage = lazy(() =>
  import('@/app/dashboard/documents/page').then(mod => ({
    default: mod.default
  }))
);

export const OrganizationsPage = lazy(() =>
  import('@/app/dashboard/organizations/page').then(mod => ({
    default: mod.default
  }))
);

export const UsersPage = lazy(() =>
  import('@/app/dashboard/users/page').then(mod => ({
    default: mod.default
  }))
);

export const DocumentDetailPage = lazy(() =>
  import('@/app/dashboard/documents/[id]/page').then(mod => ({
    default: mod.default
  }))
);

export const OrganizationDetailPage = lazy(() =>
  import('@/app/dashboard/organizations/[id]/page').then(mod => ({
    default: mod.default
  }))
);

export const UserDetailPage = lazy(() =>
  import('@/app/dashboard/users/[id]/page').then(mod => ({
    default: mod.default
  }))
);

export const DocumentCreatePage = lazy(() =>
  import('@/app/dashboard/documents/create/page').then(mod => ({
    default: mod.default
  }))
);

export const OrganizationCreatePage = lazy(() =>
  import('@/app/dashboard/organizations/create/page').then(mod => ({
    default: mod.default
  }))
);
