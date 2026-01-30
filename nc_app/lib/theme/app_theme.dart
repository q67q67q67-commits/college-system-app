import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../config/api_config.dart';

ThemeData get appTheme => ThemeData(
      useMaterial3: true,
      textTheme: GoogleFonts.montserratTextTheme(),
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(NarxozColors.red),
        primary: const Color(NarxozColors.red),
        secondary: const Color(NarxozColors.dark),
        brightness: Brightness.light,
      ),
      appBarTheme: const AppBarTheme(
        backgroundColor: Color(NarxozColors.dark),
        foregroundColor: Colors.white,
      ),
      cardTheme: CardThemeData(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(2)),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(2)),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(2)),
      ),
    );
