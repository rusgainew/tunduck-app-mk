/**
 * Next.js Compatible Dynamic Page Imports
 * For use with App Router dynamic imports
 */

import dynamic from 'next/dynamic';
import { LoadingWrapper } from '@/lib/loading-wrapper';

// Dashboard pages with lazy loading
export const DynamicDocumentsPage = dynamic(
  () => import('@/app/dashboard/documents/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true, // Enable SSR for documents page
  }
);

export const DynamicOrganizationsPage = dynamic(
  () => import('@/app/dashboard/organizations/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);

export const DynamicUsersPage = dynamic(
  () => import('@/app/dashboard/users/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);

// Detail pages
export const DynamicDocumentDetailPage = dynamic(
  () => import('@/app/dashboard/documents/[id]/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);

export const DynamicOrganizationDetailPage = dynamic(
  () => import('@/app/dashboard/organizations/[id]/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);

export const DynamicUserDetailPage = dynamic(
  () => import('@/app/dashboard/users/[id]/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);

// Create pages
export const DynamicDocumentCreatePage = dynamic(
  () => import('@/app/dashboard/documents/create/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);

export const DynamicOrganizationCreatePage = dynamic(
  () => import('@/app/dashboard/organizations/create/page').then(mod => mod.default),
  {
    loading: LoadingWrapper,
    ssr: true,
  }
);
