import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import MenuIcon from "../../../assets/icons/menu.svg";
import LeaderboardIcon from "../../../assets/icons/leaderboard.svg";
import TorrentsIcon from "../../../assets/icons/torrent.svg";
import ProfileIcon from "../../../assets/icons/profile.svg";
import { useRouter } from "../../../routes/hooks/use-pathname";

const navItems = [
  { label: "Leaderboard", icon: LeaderboardIcon, path: "/leader" },
  { label: "Torrents", icon: TorrentsIcon, path: "/torrents" },
  { label: "Profile", icon: ProfileIcon, path: "/profile" },
];

export default function NavDropdown() {
  const [isOpen, setIsOpen] = useState(false);
  const router = useRouter();

  return (
    <nav className="fixed top-4 left-5 transform z-50">
      <motion.div
        className="flex items-center justify-center w-fit px-4 py-2 rounded-xl border-2 border-[#F90627] bg-zinc-900/50 backdrop-blur-sm shadow-lg cursor-pointer"
        onHoverStart={() => setIsOpen(true)}
        onHoverEnd={() => setIsOpen(false)}
      >
        <img src={MenuIcon} alt="Menu" className="w-8 h-8" />
        <AnimatePresence>
          {isOpen && (
            <motion.div
              className="flex ml-4 gap-6"
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -10, transition: { ease: "easeInOut" } }}
            >
              {navItems.map((item, index) => (
                <motion.div
                  key={index}
                  className="flex flex-col items-center cursor-pointer"
                  whileHover={{ scale: 1.1 }}
                  onClick={() => router.push(item.path)}
                >
                  <img
                    src={item.icon}
                    alt={item.label}
                    className="w-12 h-12 drop-shadow-[0_0_15px_#F90627]"
                  />
                  <span className="mt-1 text-white font-nunito text-sm select-none">
                    {item.label}
                  </span>
                </motion.div>
              ))}
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </nav>
  );
}
