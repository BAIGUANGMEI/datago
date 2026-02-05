import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig, i18n} = useDocusaurusContext();
  const isEnglish = i18n.currentLocale === 'en';
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <p className={styles.heroDescription}>
          {isEnglish
            ? 'DataFrame, Series, GroupBy, Merge/Join, Parallel Processing - all in Go'
            : 'DataFrame、Series、分组聚合、数据合并、并行处理 - Go 语言原生实现'}
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/intro">
            {isEnglish ? 'Get Started' : '快速开始'}
          </Link>
          <Link
            className="button button--outline button--lg"
            to="/docs/examples"
            style={{marginLeft: '1rem', color: 'white', borderColor: 'white'}}>
            {isEnglish ? 'View Examples' : '查看示例'}
          </Link>
        </div>
      </div>
    </header>
  );
}

function PerformanceSection() {
  const {i18n} = useDocusaurusContext();
  const isEnglish = i18n.currentLocale === 'en';
  
  return (
    <section className={styles.performanceSection}>
      <div className="container">
        <Heading as="h2" className={styles.sectionTitle}>
          {isEnglish ? 'Performance' : '性能表现'}
        </Heading>
        <p className={styles.sectionSubtitle}>
          {isEnglish
            ? '2x faster than pandas for Excel operations'
            : 'Excel 读取速度是 pandas 的 2 倍'}
        </p>
        <div className={styles.benchmarkTable}>
          <table>
            <thead>
              <tr>
                <th>{isEnglish ? 'Dataset' : '数据集'}</th>
                <th>DataGo</th>
                <th>pandas</th>
                <th>polars</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>15K rows × 11 cols</td>
                <td className={styles.highlight}>0.21s</td>
                <td>0.51s</td>
                <td>0.10s</td>
              </tr>
              <tr>
                <td>271K rows × 16 cols</td>
                <td className={styles.highlight}>5.76s</td>
                <td>11.01s</td>
                <td>2.18s</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </section>
  );
}

function CodePreview() {
  const {i18n} = useDocusaurusContext();
  const isEnglish = i18n.currentLocale === 'en';
  
  const codeExample = `// Create DataFrame
df, _ := dataframe.New(map[string][]interface{}{
    "product": {"A", "B", "A", "B"},
    "sales":   {100.0, 150.0, 200.0, 120.0},
})

// GroupBy aggregation
gb, _ := df.GroupBy("product")
stats := gb.Sum("sales")

// Merge DataFrames
result, _ := dataframe.Merge(left, right, MergeOptions{
    How: InnerJoin,
    On:  []string{"id"},
})

// Parallel processing
result := df.ParallelFilter(func(row Row) bool {
    return row.Get("sales").(float64) > 100
})`;

  return (
    <section className={styles.codeSection}>
      <div className="container">
        <Heading as="h2" className={styles.sectionTitle}>
          {isEnglish ? 'Simple & Powerful' : '简洁而强大'}
        </Heading>
        <p className={styles.sectionSubtitle}>
          {isEnglish
            ? 'Familiar pandas-like API with Go performance'
            : '熟悉的 pandas 风格 API，Go 原生性能'}
        </p>
        <div className={styles.codePreview}>
          <pre><code>{codeExample}</code></pre>
        </div>
      </div>
    </section>
  );
}

export default function Home(): ReactNode {
  const {siteConfig, i18n} = useDocusaurusContext();
  const isEnglish = i18n.currentLocale === 'en';
  return (
    <Layout
      title={isEnglish ? 'DataGo - Go Data Analysis Library' : 'DataGo - Go 数据分析库'}
      description={
        isEnglish
          ? 'DataGo: High-performance data analysis library for Go with DataFrame, GroupBy, Merge/Join, and parallel processing'
          : 'DataGo：高性能 Go 数据分析库，支持 DataFrame、分组聚合、数据合并、并行处理'
      }>
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <PerformanceSection />
        <CodePreview />
      </main>
    </Layout>
  );
}
