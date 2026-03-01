const path = require('path');
const childProcess = require('child_process');

function run(command) {
  try {
    return childProcess.execSync(command, { stdio: ['ignore', 'pipe', 'ignore'] }).toString().trim();
  } catch (e) {
    return '';
  }
}

function parseRepositoryUrl(remote) {
  if (!remote) {
    return '';
  }

  // git@github.com:owner/repo.git
  const ssh = remote.match(/^git@github\.com:([^/]+)\/(.+?)(?:\.git)?$/i);
  if (ssh) {
    return `https://github.com/${ssh[1]}/${ssh[2]}`;
  }

  // https://github.com/owner/repo.git
  const https = remote.match(/^https?:\/\/github\.com\/([^/]+)\/(.+?)(?:\.git)?$/i);
  if (https) {
    return `https://github.com/${https[1]}/${https[2]}`;
  }

  return '';
}

function repositoryUrlFromActionsEnv() {
  // GitHub Actions always provides owner/repo for the running repository.
  if (process.env.GITHUB_REPOSITORY) {
    const serverUrl = process.env.GITHUB_SERVER_URL || 'https://github.com';
    return `${serverUrl}/${process.env.GITHUB_REPOSITORY}`;
  }

  return '';
}

function getDefaultBranch() {
  // Prefer explicit env from workflow contexts when available.
  if (process.env.DOCS_DEFAULT_BRANCH) {
    return process.env.DOCS_DEFAULT_BRANCH;
  }

  if (process.env.GITHUB_REF_NAME) {
    return process.env.GITHUB_REF_NAME;
  }

  const ref = run('git symbolic-ref --short refs/remotes/origin/HEAD');
  if (!ref) {
    return 'master';
  }

  const parts = ref.split('/');
  return parts[parts.length - 1] || 'master';
}

const originRemote = run('git config --get remote.origin.url');
const repositoryUrl = repositoryUrlFromActionsEnv() || parseRepositoryUrl(originRemote) || 'https://github.com/jandedobbeleer/aliae';
const repositoryParts = repositoryUrl.split('/');
const repositoryOwner = repositoryParts[3] || 'jandedobbeleer';
const repositoryName = repositoryParts[4] || 'aliae';
const repositorySlug = `${repositoryOwner}/${repositoryName}`;
const defaultBranch = getDefaultBranch();

module.exports = {
  title: 'aliae',
  tagline: 'Cross platform shell management.',
  url: process.env.DOCS_URL || 'https://aliae.dev',
  baseUrl: process.env.DOCS_BASE_URL || '/',
  markdown: {
    mermaid: true,
  },
  favicon: 'img/favicon.ico',
  organizationName: repositoryOwner,
  projectName: repositoryName,
  customFields: {
    repositoryUrl,
    repositoryOwner,
    repositoryName,
    repositorySlug,
  },
  onBrokenLinks: 'ignore',
  plugins: [
    path.resolve(__dirname, 'plugins', 'appinsights'),
    'docusaurus-node-polyfills'
  ],
  themes: ['@docusaurus/theme-mermaid'],
  stylesheets: [
    "https://rsms.me/inter/inter.css",
    "https://fonts.googleapis.com/css2?family=Fira+Code&display=swap"
  ],
  themeConfig: {
    prism: {
      additionalLanguages: ['powershell', 'lua', 'jsstacktrace', 'yaml'],
    },
    navbar: {
      title: 'aliae 🌱',
      items: [
        {
          to: 'docs',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {
          href: 'https://github.com/sponsors/JanDeDobbeleer',
          label: 'Sponsor',
          position: 'left',
        },
        {
          href: repositoryUrl,
          className: 'header-github-link',
          'aria-label': 'GitHub repository',
          position: 'right',
        },
        {
          href: 'https://www.gitkraken.com/invite/nQmDPR9D',
          className: 'header-gk-link',
          'aria-label': 'GitKraken',
          position: 'right',
        },
        {
          href: 'https://discord.gg/n7E3DkXssv',
          className: 'header-discord-link',
          'aria-label': 'Discord',
          position: 'right',
        },
        {
          href: 'https://staging.bsky.app/profile/aliae.dev',
          className: 'header-bluesky-link',
          'aria-label': 'Bluesky',
          position: 'right',
        }
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'How to',
          items: [
            {
              label: 'Getting started',
              to: 'docs/',
            },
            {
              label: 'Contributing',
              to: 'docs/contributing/started',
            },
          ],
        },
        {
          title: 'Social',
          items: [
            {
              label: 'GitHub',
              href: repositoryUrl,
            },
            {
              label: 'Discord',
              href: 'https://discord.gg/n7E3DkXssv',
            },
            {
              label: 'Bluesky',
              href: 'https://staging.bsky.app/profile/aliae.dev',
            }
          ],
        },
        {
          title: 'Links',
          items: [
            {
              label: 'Sponsor',
              href: 'https://github.com/sponsors/JanDeDobbeleer',
            },
            {
              label: 'GitKraken',
              href: 'https://www.gitkraken.com/invite/nQmDPR9D',
            },
            {
              label: 'Docusaurus',
              href: 'https://github.com/facebook/docusaurus',
            },
            {
              label: 'Privacy',
              href: '/privacy',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} <a href='https://github.com/sponsors/JanDeDobbeleer' target='_blank'>Jan De Dobbeleer</a> and <a href='/docs/contributors'>contributors</a>.`,
    },
    appInsights: {
      instrumentationKey: 'b204d2e6-9a10-473c-9655-2ad79af82e5c',
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: `${repositoryUrl}/edit/${defaultBranch}/website/`,
        },
        theme: {
          customCss: [
            require.resolve('./src/css/prism-rose-pine-moon.css'),
            require.resolve('./src/css/custom.css')
          ],
        },
      },
    ],
  ],
};
