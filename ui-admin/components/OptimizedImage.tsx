/**
 * Image Optimization Component
 * Wrapper around Next.js Image for consistent optimization
 */

import Image from 'next/image';
import { CSSProperties } from 'react';

interface OptimizedImageProps {
  src: string;
  alt: string;
  width?: number;
  height?: number;
  className?: string;
  style?: CSSProperties;
  priority?: boolean;
  lazy?: boolean;
}

/**
 * Optimized Image Component
 * Automatically handles:
 * - Responsive sizing
 * - Lazy loading (by default)
 * - Format optimization (WebP with fallback)
 * - Blur placeholder (optional)
 */
export const OptimizedImage = ({
  src,
  alt,
  width = 800,
  height = 600,
  className = '',
  style,
  priority = false,
  lazy = true,
}: OptimizedImageProps) => {
  return (
    <Image
      src={src}
      alt={alt}
      width={width}
      height={height}
      className={className}
      style={style}
      priority={priority}
      loading={lazy && !priority ? 'lazy' : 'eager'}
      sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
      quality={75}
    />
  );
};

/**
 * Avatar Image Component
 * For user/org avatars with circular styling
 */
export const AvatarImage = ({
  src,
  alt,
  className = 'w-10 h-10',
}: {
  src: string;
  alt: string;
  className?: string;
}) => (
  <OptimizedImage
    src={src}
    alt={alt}
    width={40}
    height={40}
    className={`rounded-full object-cover ${className}`}
    lazy
  />
);

/**
 * Hero Image Component
 * For large hero/banner images
 */
export const HeroImage = ({
  src,
  alt,
  className = 'w-full h-96',
}: {
  src: string;
  alt: string;
  className?: string;
}) => (
  <OptimizedImage
    src={src}
    alt={alt}
    width={1200}
    height={400}
    className={`object-cover ${className}`}
    priority
  />
);

/**
 * Icon Image Component
 * For small icons with consistent sizing
 */
export const IconImage = ({
  src,
  alt,
  className = 'w-6 h-6',
}: {
  src: string;
  alt: string;
  className?: string;
}) => (
  <OptimizedImage
    src={src}
    alt={alt}
    width={24}
    height={24}
    className={className}
    lazy
  />
);
