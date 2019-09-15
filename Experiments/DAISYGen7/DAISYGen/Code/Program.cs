// Main file for DAISYGen - Convertor from audiotracks to DAISY-compatible audiobook format.
// This project was built via MonoDevelop with VS 2012 format support.

using System;
using System.Collections.Generic;
using System.Windows.Forms;
using System.Drawing;
using Gtk;

// Windows Media Player library, represented as .NET assembly
// with help of Windows Forms ActiveX Control Importer (Aximp.exe)
using AxWMPLib;

namespace DAISYGen
{	
	internal sealed class MainClass
	{
		//Global objects:
		//-- Windows Media Player hidden component
		public static AxWindowsMediaPlayer wmp;
		//-- Time comparer object for sorting time positions in an audiobook
		public static TimeComparer time_comparer = new TimeComparer ();
		//-- Path to executable
		public static string path;
		
		// Initialize hidden form with media player
		static Form wmpInit()
		{
			//Create form and corresponding control
			Form wmp_hidden_form = new Form ();
			wmp = new AxWindowsMediaPlayer ();
			wmp_hidden_form.Controls.Add (wmp);
			//Hide the form and store default volume settings
			wmp_hidden_form.FormBorderStyle = FormBorderStyle.None;
			wmp_hidden_form.BackColor = wmp_hidden_form.TransparencyKey = Color.Magenta;
			wmp_hidden_form.AllowTransparency = true;
			wmp_hidden_form.ShowInTaskbar = false;
			wmp_hidden_form.Show ();
			wmp_hidden_form.Visible = false;
			wmp.settings.volume = 100;
			wmp_hidden_form.Hide ();
			//Return hidden form
			return wmp_hidden_form;
		}

		//Main, single-threaded, entry point procedure
		[STAThread]
		public static void Main (string[] args)
		{
			//Initialize the path to executable
			path = System.Windows.Forms.Application.ExecutablePath;
			//Initialize the Windows Media Player component
			Form wmp_hidden_form = wmpInit ();
			//This program uses GTK# 2.0
			Gtk.Application.Init ();
			MainWindow win = new MainWindow ();
			//Control the status of the Windows Media Player
			wmp.StatusChange += new EventHandler (win.AxWindowsMediaPlayer1StatusChange);
			//Show the main window and maximize it 
			win.Show ();
			win.Maximize ();
			Gtk.Application.Run ();
			//Dispose of the Windows Media Player at the end 
			wmp_hidden_form.Dispose ();
		}
	}
}
