import 'package:flutter/material.dart';
import '../config/api_config.dart';

ThemeData get appTheme => ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: Color(NarxozColors.red),
        primary: Color(NarxozColors.red),
        secondary: Color(NarxozColors.dark),
        brightness: Brightness.light,
      ),
      appBarTheme: AppBarTheme(
        backgroundColor: Color(NarxozColors.dark),
        foregroundColor: Colors.white,
      ),
    );
