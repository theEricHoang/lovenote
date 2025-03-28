import { cn } from "app/lib/utils"; // Utility function for conditional classes
import type { ReactNode } from "react";
import { Link, type LinkProps } from "react-router";

interface NavButtonProps extends LinkProps {
  variant?: "primary" | "outline" | "danger";
  size?: "sm" | "md" | "lg";
  isLoading?: boolean;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
}

export default function NavButton({
  variant = "primary",
  size = "md",
  leftIcon,
  rightIcon,
  className,
  children,
  ...props
}: NavButtonProps) {
  const baseStyles =
    "flex items-center justify-center rounded-4xl font-medium transition-all focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed";

  const variants = {
    primary: "bg-rose-400 text-white hover:bg-rose-600",
    outline: "border border-rose-400 text-rose-500 hover:bg-gray-100",
    danger: "bg-red-600 text-white hover:bg-red-700",
  };

  const sizes = {
    sm: "px-3 py-1 text-sm",
    md: "px-4 py-2 text-base",
    lg: "px-5 py-3 text-lg",
  };

  return (
    <Link
      className={cn(
        baseStyles,
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {leftIcon && <span className="mr-2">{leftIcon}</span>}
      {children}
      {rightIcon && <span className="ml-2">{rightIcon}</span>}
    </Link>
  );
}
