export const COMPANY = {
  NAME: "GIT_WIZ",
  SHORT_NAME: "GIT_WIZ",
  MARKETING_URL: "https://gitwiz.dev",
  DOCS_URL: "https://gitwiz.dev/docs",
  GITHUB_URL: "https://github.com/sriram687/GIT_WIZ",
  GITHUB_REPO: "sriram687/GIT_WIZ",
  STATUS_URL: "https://status.gitwiz.dev",
  TRUST_URL: "https://gitwiz.dev/privacy",
  MAIL_TO: "mailto:hello@gitwiz.dev",
  X_URL: "https://twitter.com/gitwiz",
  LINKEDIN_URL: "https://linkedin.com/company/gitwiz",
  YOUTUBE_URL: "https://youtube.com/@gitwiz",
  DISCORD_URL: "https://discord.com/invite/UZv7JjxbwG",
  EMAIL_DOMAIN: "@gitwiz.dev",
  FOUNDERS_MAIL_TO: "mailto:founders@gitwiz.dev",
  FOUNDERS_EMAIL: "founders@gitwiz.dev",
  REPORT_ISSUE_URL: "https://github.com/sriram687/GIT_WIZ/issues/new",
  LICENSE: "Apache-2.0",
  LICENSE_URL: "https://github.com/sriram687/GIT_WIZ/blob/main/LICENSE",
} as const;

export const THEME_STORAGE_KEY = "gitwiz-theme";
export const POSTHOG_COOKIE_NAME = "ph_phc_";

export const OPEN_ROLES = [] as { title: string; url: string; location: string }[];

export const PLATFORMS = {
  MACOS: "macos",
  WINDOWS: "windows",
  LINUX: "linux",
} as const;

export const GITHUB_STARS_URL = "https://api.github.com/repos/sriram687/GIT_WIZ";

export const DOWNLOAD_URL_MAC_ARM64 = "https://github.com/sriram687/GIT_WIZ/releases/latest/download/git-wiz-darwin-arm64.zip";
export const DOWNLOAD_URL_MAC_X64 = "https://github.com/sriram687/GIT_WIZ/releases/latest/download/git-wiz-darwin-x64.zip";
export const DOWNLOAD_URL_WINDOWS = "https://github.com/sriram687/GIT_WIZ/releases/latest/download/git-wiz-win32-x64.exe";
export const DOWNLOAD_URL_LINUX = "https://github.com/sriram687/GIT_WIZ/releases/latest/download/git-wiz-linux-x64.AppImage";

export const AGENT_HARNESSES = 23;
export const TAGLINE = "Supervise parallel AI coding agents in isolated Git worktrees.";
export const HERO_SUBHEADLINE = "Run a fleet of coding agents with zero workspace collisions, live terminal oversight, and automated CI feedback loops.";
export const HERO_SECONDARY_SUBHEADLINE = "Isolated Git worktree sandboxes for Claude Code, Cursor, Aider, and any CLI agent. Review every change from one local supervisor dashboard. Free and open source.";

export const NAV_ITEMS = [
  { label: "Demo", href: "/#see-it" },
  { label: "Features", href: "/#features" },
  { label: "Changelog", href: "/changelog" },
] as const;
