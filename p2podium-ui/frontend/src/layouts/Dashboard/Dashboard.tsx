import { usePathname, useRouter } from "../../routes/hooks/use-pathname";
import Left from "../../assets/icons/chevron-left.svg";
import { motion } from "framer-motion";
import NavBar from "./components/NavBar";

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const router = useRouter();
  const currentPath = usePathname();
  const isHomePage = currentPath === "/";

  const buttonVariants = {
    rest: { scale: 1, rotate: 0 },
    hover: { scale: 1.1, rotate: 10 },
    tap: { scale: 0.95, rotate: -5 },
  };

  return (
    <article className="h-screen w-screen flex items-center justify-center relative">
      {!isHomePage && (
        <motion.button
          onClick={router.back}
          variants={buttonVariants}
          initial="rest"
          whileHover="hover"
          whileTap="tap"
          animate="rest"
          className="fixed top-5 left-5 h-[50px] w-[50px] cursor-pointer flex items-center justify-center rounded-lg bg-zinc-900/50 border border-zinc-800 shadow-lg"
        >
          <img
            src={Left}
            className="h-10 w-10 pointer-events-none"
            alt="Go back"
          />
        </motion.button>
      )}
      <NavBar/>

      {children}
    </article>
  );
}

