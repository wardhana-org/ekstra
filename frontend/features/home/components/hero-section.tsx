import styles from "../styles/hero-section.module.scss"

export default function HeroSection() {
    return(
        <section className={styles.section}>
            <div className={styles.header}>
                <div className={styles.logoContainerBackground}>
                    <div className={styles.logoContainer}>
                        <span>logo</span>
                    </div>
                </div>
                <div className={styles.navContainer}>
                    <nav>
                        <a href="">test</a>
                    </nav>
                </div>
                <div className={styles.authContainerBackground}>
                    <div className={styles.authContainer}>
                        <nav>
                            <a href="">Login</a>
                        </nav>
                </div>
                </div>
           </div>
           <div className={styles.heroContent}>
                <p>test</p>
                <div className={styles.modeContainer}>
                    <span>test</span>
                </div>
           </div>
       </section>
    )
}
