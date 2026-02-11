import 'package:flutter/material.dart';

/// Тема приложения: палитра collegenarxoz.kz, минимализм, минимальные закругления.
ThemeData buildNcTheme() {
  const primary = Color(0xFFD7263D);
  const background = Color(0xFFF5F5F7);
  const surface = Colors.white;
  const textColor = Color(0xFF222222);

  final base = ThemeData.light();

  return base.copyWith(
    colorScheme: ColorScheme.fromSeed(
      seedColor: primary,
      primary: primary,
      secondary: Colors.grey.shade700,
      surface: surface,
    ),
    scaffoldBackgroundColor: background,
    appBarTheme: const AppBarTheme(
      backgroundColor: surface,
      elevation: 0,
      foregroundColor: textColor,
      centerTitle: true,
    ),
    cardTheme: CardThemeData(
      color: surface,
      elevation: 1,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
    ),
    bottomNavigationBarTheme: const BottomNavigationBarThemeData(
      selectedItemColor: primary,
      unselectedItemColor: Colors.grey,
      type: BottomNavigationBarType.fixed,
    ),
    navigationRailTheme: NavigationRailThemeData(
      selectedIconTheme: const IconThemeData(color: primary),
      unselectedIconTheme: IconThemeData(color: Colors.grey.shade600),
    ),
    textTheme: base.textTheme.copyWith(
      headlineSmall: const TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: textColor,
      ),
      bodyMedium: const TextStyle(
        fontSize: 14,
        color: textColor,
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(6),
        borderSide: const BorderSide(color: Colors.grey),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(6),
        borderSide: const BorderSide(color: primary),
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: primary,
        foregroundColor: Colors.white,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(6),
        ),
        minimumSize: const Size(double.infinity, 44),
      ),
    ),
  );
}
