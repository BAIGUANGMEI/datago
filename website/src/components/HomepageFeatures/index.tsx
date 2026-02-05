import type {ReactNode} from 'react';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  description: ReactNode;
  link?: string;
  linkText?: string;
};

const buildFeatureList = (isEnglish: boolean): FeatureItem[] => [
  {
    title: isEnglish ? 'DataFrame & Series' : 'DataFrame 与 Series',
    description: (
      <>
        {isEnglish
          ? 'Familiar pandas-like API with DataFrame and Series. Filter, sort, select, and transform data with ease.'
          : '熟悉的 pandas 风格 API，支持 DataFrame 和 Series。轻松实现筛选、排序、选择和转换。'}
      </>
    ),
    link: '/docs/dataframe',
    linkText: isEnglish ? 'Learn more' : '了解更多',
  },
  {
    title: isEnglish ? 'GroupBy Aggregation' : 'GroupBy 分组聚合',
    description: (
      <>
        {isEnglish
          ? 'Powerful grouping and aggregation operations. Sum, Mean, Count, and custom aggregations.'
          : '强大的分组聚合功能。支持 Sum、Mean、Count 等内置聚合及自定义聚合。'}
      </>
    ),
    link: '/docs/groupby',
    linkText: isEnglish ? 'Learn more' : '了解更多',
  },
  {
    title: isEnglish ? 'Merge & Join' : '数据合并 Merge/Join',
    description: (
      <>
        {isEnglish
          ? 'SQL-like join operations. Inner, Left, Right, and Outer joins with multi-key support.'
          : 'SQL 风格的表关联操作。支持 Inner、Left、Right、Outer Join 及多键合并。'}
      </>
    ),
    link: '/docs/merge',
    linkText: isEnglish ? 'Learn more' : '了解更多',
  },
  {
    title: isEnglish ? 'Parallel Processing' : '并行处理',
    description: (
      <>
        {isEnglish
          ? 'Leverage Go concurrency for big data. ParallelApply, ParallelFilter, and parallel aggregations.'
          : '充分利用 Go 并发优势处理大数据。支持并行 Apply、Filter 和聚合操作。'}
      </>
    ),
    link: '/docs/parallel',
    linkText: isEnglish ? 'Learn more' : '了解更多',
  },
  {
    title: isEnglish ? 'Excel & CSV I/O' : 'Excel 与 CSV 读写',
    description: (
      <>
        {isEnglish
          ? 'High-performance file I/O. 2x faster than pandas for Excel operations.'
          : '高性能文件读写。Excel 读取速度是 pandas 的 2 倍。'}
      </>
    ),
    link: '/docs/io-excel',
    linkText: isEnglish ? 'Learn more' : '了解更多',
  },
  {
    title: isEnglish ? 'Type Safe' : '类型安全',
    description: (
      <>
        {isEnglish
          ? "Go's static typing provides compile-time checks. Multiple data types with automatic inference."
          : 'Go 静态类型提供编译时检查。支持多种数据类型及自动推断。'}
      </>
    ),
    link: '/docs/series',
    linkText: isEnglish ? 'Learn more' : '了解更多',
  },
];

function Feature({title, description, link, linkText}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className={clsx('text--center padding-horiz--md', styles.featureCard)}>
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
        {link && (
          <Link className={styles.featureLink} to={link}>
            {linkText} →
          </Link>
        )}
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
            {isEnglish ? 'Everything you need for data analysis in Go' : '为 Go 打造的全能数据分析工具'}
          </Heading>
          <p>
            {isEnglish
              ? 'High-performance, pandas-like API, fully leveraging Go concurrency'
              : '高性能、类 pandas API、充分利用 Go 并发优势'}
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
