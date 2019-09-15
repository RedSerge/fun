/*
 * Created by SharpDevelop.
 * User: RedSerge
 * Date: 5/25/2019
 * Time: 4:30 PM
 * 
 * To change this template use Tools | Options | Coding | Edit Standard Headers.
 */
using System;
using System.Collections.Generic;
using System.Drawing;
using System.Windows.Forms;
using System.Timers;

namespace MinBrowser
{
	/// <summary>
	/// Description of MainForm.
	/// </summary>
	public partial class MainForm : Form
	{
		public MainForm()
		{
			//
			// The InitializeComponent() call is required for Windows Forms designer support.
			//
			InitializeComponent();
			
			//
			// TODO: Add constructor code after the InitializeComponent() call.
			//
		}
		void Button1Click(object sender, EventArgs e)
		{
			string query="Pterocles burchelli";
			string new_q=query.Replace(' ','+');
			string s="https://www.wolframalpha.com/input/?i="+new_q;
			webBrowser1.Navigate(s);			
		}
		void WebBrowser1DocumentCompleted(object sender, WebBrowserDocumentCompletedEventArgs e)
		{			
			
		}
		void Button2Click(object sender, EventArgs e)
		{
			System.IO.File.WriteAllText("fetch.txt", webBrowser1.Document.Body.InnerHtml);			
		}
		int inum=0;
		void Button3Click(object sender, EventArgs e)
		{
			textBox2.Lines=System.IO.File.ReadAllLines("available.txt");			
			inum=0;
			string query=textBox2.Lines[0];
			string new_q=query.Replace(' ','+');
			string s="https://www.wolframalpha.com/input/?i="+new_q;
			timer1.Start();
			webBrowser1.Navigate(s);			
		}
		void Timer1Tick(object sender, EventArgs e)
		{
			if (inum>=textBox2.Lines.Length) {
				timer1.Stop();
			} else {				
				System.IO.File.WriteAllText(textBox2.Lines[inum]+".wiki.txt", webBrowser1.Document.Body.InnerHtml);
				inum+=1;
				string query=textBox2.Lines[inum];
				string new_q=query.Replace(' ','+');
				string s="https://www.wolframalpha.com/input/?i="+new_q;
				webBrowser1.Navigate(s);
			}
		}
	}
}
