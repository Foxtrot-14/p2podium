import { motion } from "framer-motion";

export default function Home() {
  return (
    <article className="h-fit w-fit flex flex-col items-center justify-center gap-6">
      <motion.h1
        key="title"
        initial={{ opacity: 0, y: 50 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{
          duration: 1,
          ease: "easeInOut",
        }}
        className="text-6xl text-[#F90627] tracking-widest select-none font-orbitron"
      >
        P2Podium
      </motion.h1>
   </article>
  );
}

