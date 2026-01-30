import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
};

const buildFeatureList = (isEnglish: boolean): FeatureItem[] => [
  {
    title: isEnglish ? 'Pandas-like API' : '类似 pandas 的 API',
    description: (
      <>
        {isEnglish
          ? 'Familiar DataFrame/Series structures and common data ops.'
          : '熟悉的 DataFrame / Series 结构与常用数据操作，快速上手。'}
      </>
    ),
  },
  {
    title: isEnglish ? 'Columnar Performance' : '高效的列式处理',
    description: (
      <>
        {isEnglish
          ? 'Column-first design for analytics and batch computation.'
          : '以列为中心的数据结构，适用于分析场景与批量计算。'}
      </>
    ),
  },
  {
    title: isEnglish ? 'I/O for Common Formats' : '多种类型读写支持',
    description: (
      <>
        {isEnglish
          ? 'Built-in I/O helpers for everyday data workflows.'
          : '内置多种类型读写接口，方便与日常数据流程对接。'}
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className={clsx('text--center padding-horiz--md', styles.featureCard)}>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): ReactNode {
  const {i18n} = useDocusaurusContext();
  const isEnglish = i18n.currentLocale === 'en';
  const featureList = buildFeatureList(isEnglish);

  return (
    <section className={styles.features}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <Heading as="h2">
            {isEnglish ? 'Data analysis for Go' : '为 Go 而生的数据分析库'}
          </Heading>
          <p>
            {isEnglish
              ? 'Lightweight, intuitive, and focused on structured data & Excel workflows.'
              : '轻量、直观、可扩展，专注于结构化数据与 Excel 工作流。'}
          </p>
        </div>
        <div className="row">
          {featureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
