import { motion } from "framer-motion";

export default function Home() {
  return (
    <article className="h-fit w-fit flex flex-col items-center justify-center gap-6">
      <motion.h1
        key="title"
        initial={{ opacity: 0, y: 50 }}
        animate={{
          opacity: [0, 1, 1, 0],
          y: [50, 0, 0, -20],
        }}
        transition={{
          duration: 4,
          times: [0, 0.25, 0.75, 1],
          ease: "easeInOut",
        }}
        className="text-6xl text-[#F90627] tracking-widest select-none font-orbitron"
      >
        P2Podium
      </motion.h1>
    </article>
  );
}
