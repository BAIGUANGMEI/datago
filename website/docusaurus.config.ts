import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'DataGo',
  tagline: '高性能 Go 数据分析库 | DataFrame, GroupBy, Merge, 并行处理',
  // favicon: 'img/logo-datago.png',

  // Future flags, see https://docusaurus.io/docs/api/docusaurus-config#future
  future: {
    v4: true, // Improve compatibility with the upcoming Docusaurus v4
  },

  // Set the production url of your site here
  url: 'https://baiguangmei.github.io',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/datago/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'BAIGUANGMEI', // Usually your GitHub org/user name.
  projectName: 'datago', // Usually your repo name.
  trailingSlash: false,

  onBrokenLinks: 'throw',

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'zh-Hans',
    locales: ['zh-Hans'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts'
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    image: 'img/logo-datago.png',
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'DataGo',
      // logo: {
      //   alt: 'DataGo Logo',
      //   src: 'img/logo-datago.svg',
      // },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: '文档',
        },
        {
          to: '/docs/examples',
          position: 'left',
          label: '示例',
        },
        // {
        //   type: 'localeDropdown',
        //   position: 'right',
        // },
        {
          href: 'https://github.com/BAIGUANGMEI/datago',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: '文档',
          items: [
            {
              label: '快速开始',
              to: '/docs/intro',
            },
            {
              label: 'DataFrame',
              to: '/docs/dataframe',
            },
            {
              label: 'GroupBy',
              to: '/docs/groupby',
            },
            {
              label: 'Merge/Join',
              to: '/docs/merge',
            },
          ],
        },
        {
          title: '功能',
          items: [
            {
              label: '并行处理',
              to: '/docs/parallel',
            },
            {
              label: 'Excel 读写',
              to: '/docs/io-excel',
            },
            {
              label: 'CSV 读写',
              to: '/docs/io-csv',
            },
            {
              label: '示例',
              to: '/docs/examples',
            },
          ],
        },
        {
          title: '社区',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/BAIGUANGMEI/datago',
            },
            {
              label: 'Issues',
              href: 'https://github.com/BAIGUANGMEI/datago/issues',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} DataGo. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['go', 'bash'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
