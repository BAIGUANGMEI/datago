import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.
 */
const sidebars: SidebarsConfig = {
  tutorialSidebar: [
    'intro',
    'getting-started',
    {
      type: 'category',
      label: '核心数据结构',
      items: ['dataframe', 'series', 'data-index'],
    },
    {
      type: 'category',
      label: '高级功能',
      items: ['groupby', 'merge', 'parallel'],
    },
    {
      type: 'category',
      label: '数据 I/O',
      items: ['io-excel', 'io-csv'],
    },
    'examples',
  ],
};

export default sidebars;
