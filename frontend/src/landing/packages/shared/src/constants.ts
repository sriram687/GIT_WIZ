export const COMPANY = {
  NAME: "GIT_WIZ",
  SHORT_NAME: "AO",
  MARKETING_URL: "https://aoagents.dev",
  DOCS_URL: "https://aoagents.dev/docs",
  GITHUB_URL: "https://github.com/sriram687/GIT_WIZ",
  GITHUB_REPO: "sriram687/GIT_WIZ",
  STATUS_URL: "https://status.aoagents.dev",
  TRUST_URL: "https://aoagents.dev/privacy",
  MAIL_TO: "mailto:hello@aoagents.dev",
  X_URL: "https://twitter.com/aoagents",
  LINKEDIN_URL: "https://linkedin.com/company/aoagents",
  YOUTUBE_URL: "https://youtube.com/@aoagents",
  DISCORD_URL: "https://discord.com/invite/UZv7JjxbwG",
  EMAIL_DOMAIN: "@aoagents.dev",
  FOUNDERS_MAIL_TO: "mailto:founders@aoagents.dev",
  FOUNDERS_EMAIL: "founders@aoagents.dev",
  REPORT_ISSUE_URL: "https://github.com/sriram687/GIT_WIZ/issues/new",
  LICENSE: "Apache-2.0",
  LICENSE_URL: "https://github.com/sriram687/GIT_WIZ/blob/main/LICENSE",
} as const;

export const THEME_STORAGE_KEY = "ao-theme";
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
export const TAGLINE = "Turn agent chaos into managable work.";
export const HERO_SUBHEADLINE = "Run a fleet of coding agents while keeping branches, reviews, and CI failures managable.";
export const HERO_SECONDARY_SUBHEADLINE = "Isolated workspaces for Claude Code, Codex, and any CLI agent. Review every change from one dashboard. Free and open source.";

export const NAV_ITEMS = [
  { label: "Demo", href: "/#see-it" },
  { label: "Features", href: "/#features" },
  { label: "Changelog", href: "/changelog" },
] as const;
